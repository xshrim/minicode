# Harbor备份恢复方案测试

## 背景

针对Harbor v1.10.6版制定离线备份恢复方案, 保证Harbor集群崩溃后能够通过备份进行数据恢复.

## 方案

可选方案:

- 通过Harbor API进行备份恢复
- 通过物理文件冷备进行备份恢复

> 备份对象可以是线上主Harbor, 也可以是开启同步机制的备Harbor(同步包括镜像和元数据同步)

API备份的方式备份和恢复过程耗时较长, 且存在对Harbor集群的压力.

文件冷备的方式备份和恢复过程耗时短, 对Harbor集群无压力.

基于以上原因, 采取物理文件冷备方案. 需对该方案的可行性和一致性进行验证.

冷备方式分为对registry, chart, redis, secret的备份恢复和对postgresql数据的备份恢复. 前者直接物理备份, 后者则可以采取逻辑备份或物理备份. 将分别验证这两种备份方式. 这里称作**全物理备份恢复**和**半物理备份恢复**.

## 环境

通过docker compose启动harbor集群, 并编写集群控制脚本.

- docker-compose.yaml

  ```yaml
  version: '2.3'
  services:
    log:
      image: goharbor/harbor-log:v1.10.6
      container_name: harbor-log
      restart: always
      dns_search: .
      cap_drop:
        - ALL
      cap_add:
        - CHOWN
        - DAC_OVERRIDE
        - SETGID
        - SETUID
      volumes:
        - /var/log/harbor/:/var/log/docker/
        - ./common/config/log/logrotate.conf:/etc/logrotate.d/logrotate.conf
        - ./common/config/log/rsyslog_docker.conf:/etc/rsyslog.d/rsyslog_docker.conf
      ports:
        - 127.0.0.1:1514:10514
      networks:
        - harbor
    registry:
      image: goharbor/registry-photon:v1.10.6
      container_name: registry
      restart: always
      cap_drop:
        - ALL
      cap_add:
        - CHOWN
        - SETGID
        - SETUID
      volumes:
        - ./data/registry:/storage
        - ./common/config/registry/:/etc/registry/
        - ./common/secret/registry/root.crt:/etc/registry/root.crt
      networks:
        - harbor
      dns_search: .
      depends_on:
        - log
      logging:
        driver: "syslog"
        options:
          syslog-address: "tcp://127.0.0.1:1514"
          tag: "registry"
    registryctl:
      image: goharbor/harbor-registryctl:v1.10.6
      container_name: registryctl
      env_file:
        - ./common/config/registryctl/env
      restart: always
      cap_drop:
        - ALL
      cap_add:
        - CHOWN
        - SETGID
        - SETUID
      volumes:
        - ./data/registry:/storage
        - ./common/config/registry/:/etc/registry/
        - ./common/config/registryctl/config.yml:/etc/registryctl/config.yml
      networks:
        - harbor
      dns_search: .
      depends_on:
        - log
      logging:
        driver: "syslog"
        options:
          syslog-address: "tcp://127.0.0.1:1514"
          tag: "registryctl"
    postgresql:
      image: goharbor/harbor-db:v1.10.6
      container_name: harbor-db
      restart: always
      cap_drop:
        - ALL
      cap_add:
        - CHOWN
        - DAC_OVERRIDE
        - SETGID
        - SETUID
      volumes:
        - ./data/database:/var/lib/postgresql/data
      networks:
        - harbor
      dns_search: .
      env_file:
        - ./common/config/db/env
      depends_on:
        - log
      logging:
        driver: "syslog"
        options:
          syslog-address: "tcp://127.0.0.1:1514"
          tag: "postgresql"
    core:
      image: goharbor/harbor-core:v1.10.6
      container_name: harbor-core
      env_file:
        - ./common/config/core/env
      restart: always
      cap_drop:
        - ALL
      cap_add:
        - SETGID
        - SETUID
      volumes:
        - ./data/ca_download/:/etc/core/ca/
        - ./data/psc/:/etc/core/token/
        - ./data/:/data/
        - ./common/config/core/certificates/:/etc/core/certificates/
        - ./common/config/core/app.conf:/etc/core/app.conf
        - ./common/secret/core/private_key.pem:/etc/core/private_key.pem
        - ./common/secret/keys/secretkey:/etc/core/key
      networks:
        - harbor
      dns_search: .
      depends_on:
        - log
        - registry
        - redis
        - postgresql
      logging:
        driver: "syslog"
        options:
          syslog-address: "tcp://127.0.0.1:1514"
          tag: "core"
    portal:
      image: goharbor/harbor-portal:v1.10.6
      container_name: harbor-portal
      restart: always
      cap_drop:
        - ALL
      cap_add:
        - CHOWN
        - SETGID
        - SETUID
        - NET_BIND_SERVICE
      networks:
        - harbor
      dns_search: .
      depends_on:
        - log
      logging:
        driver: "syslog"
        options:
          syslog-address: "tcp://127.0.0.1:1514"
          tag: "portal"
  
    jobservice:
      image: goharbor/harbor-jobservice:v1.10.6
      container_name: harbor-jobservice
      env_file:
        - ./common/config/jobservice/env
      restart: always
      cap_drop:
        - ALL
      cap_add:
        - CHOWN
        - SETGID
        - SETUID
      volumes:
        - ./data/job_logs:/var/log/jobs
        - ./common/config/jobservice/config.yml:/etc/jobservice/config.yml
      networks:
        - harbor
      dns_search: .
      depends_on:
        - core
      logging:
        driver: "syslog"
        options:
          syslog-address: "tcp://127.0.0.1:1514"
          tag: "jobservice"
    redis:
      image: goharbor/redis-photon:v1.10.6
      container_name: redis
      restart: always
      cap_drop:
        - ALL
      cap_add:
        - CHOWN
        - SETGID
        - SETUID
      volumes:
        - ./data/redis:/var/lib/redis
      networks:
        - harbor
      dns_search: .
      depends_on:
        - log
      logging:
        driver: "syslog"
        options:
          syslog-address: "tcp://127.0.0.1:1514"
          tag: "redis"
    proxy:
      image: goharbor/nginx-photon:v1.10.6
      container_name: nginx
      restart: always
      cap_drop:
        - ALL
      cap_add:
        - CHOWN
        - SETGID
        - SETUID
        - NET_BIND_SERVICE
      volumes:
        - ./common/config/nginx:/etc/nginx
      networks:
        - harbor
      dns_search: .
      ports:
        - 80:8080
      depends_on:
        - registry
        - core
        - portal
        - log
      logging:
        driver: "syslog"
        options:
          syslog-address: "tcp://127.0.0.1:1514"
          tag: "proxy"
    clair:
      container_name: clair
      image: goharbor/clair-photon:v1.10.6
      restart: always
      cap_drop:
        - ALL
      cap_add:
        - DAC_OVERRIDE
        - SETGID
        - SETUID
      cpu_quota: 50000
      dns_search: .
      depends_on:
        - log
        - postgresql
      volumes:
        - ./common/config/clair/config.yaml:/etc/clair/config.yaml
      networks:
        - harbor
      logging:
        driver: "syslog"
        options:
          syslog-address: "tcp://127.0.0.1:1514"
          tag: "clair"
      env_file:
        ./common/config/clair/clair_env
    clair-adapter:
      networks:
        - harbor
      container_name: clair-adapter
      image: goharbor/clair-adapter-photon:v1.10.6
      restart: always
      cap_drop:
        - ALL
      cap_add:
        - DAC_OVERRIDE
        - SETGID
        - SETUID
      cpu_quota: 5000http://adminset.cn0
      dns_search: .
      depends_on:
        - clair
        - redis
      logging:
        driver: "syslog"
        options:
          syslog-address: "tcp://127.0.0.1:1514"
          tag: "clair-adapter"
      env_file:
        ./common/config/clair-adapter/env
    chartmuseum:
      container_name: chartmuseum
      image: goharbor/chartmuseum-photon:v1.10.6
      restart: always
      cap_drop:
        - ALL
      cap_add:
        - CHOWN
        - DAC_OVERRIDE
        - SETGID
        - SETUID
      networks:
        - harbor
      dns_search: .
      depends_on:
        - log
      volumes:
        - ./data/chart_storage:/chart_storage
        - ./common/config/chartserver:/etc/chartserver
      logging:
        driver: "syslog"
        options:
          syslog-address: "tcp://127.0.0.1:1514"
          tag: "chartmuseum"
      env_file:
        ./common/config/chartserver/env
  networks:
    harbor:
      external: false
  ```

- 控制脚本

  ```bash
  #!/bin/bash
  
  ## filename: hbr
  
  kind="logical"
  pause="false"
  
  ##### docker-compose
  function run() {
    clean
  
    mkdir -p data/{ca_download,chart_storage,database,job_logs,psc,redis,registry,secret}
  
    chown 10000:10000 data/ca_download
    chown 10000:10000 data/chart_storage
    chown 999:999 data/database
    chown 10000:10000 data/job_logs
    chown 10000:10000 data/psc
    chown 999:999 data/redis
    chown 10000:10000 data/registry
    chown root:root data/secret
  
    up
  }
  
  function up() {
    docker-compose up -d
  }
  
  function down() {
    docker-compose down
  }
  
  function clean() {
    down
  
    docker rm -f $(docker ps -a|grep goharbor|awk '{print $1}')
    docker volume prune -f
  
    rm -rf data
  }
  
  ##### backup
  create_dir(){
    rm -rf harbor
    mkdir -p harbor/secret
    chmod 777 harbor/secret
  }
  
  create_tarball() {
    if [ -n "$1" ]; then
      backup_filename=$1
    else
      timestamp=$(date +"%Y-%m-%d-%H-%M-%S")
      backup_filename=harbor-$timestamp.tgz
    fi
    tar zcvf $backup_filename harbor
  }
  
  clean_harbor() {
    rm -rf harbor
  }
  
  export_database() {
    docker exec harbor-db sh -c 'rm -rf /home/postgres/database'
    docker exec harbor-db sh -c 'mkdir /home/postgres/database'
    docker exec harbor-db sh -c 'pg_dump -U postgres registry > /home/postgres/database/registry.bak'
    docker exec harbor-db sh -c 'pg_dump -U postgres postgres > /home/postgres/database/postgres.bak'
    docker exec harbor-db sh -c 'pg_dump -U postgres notarysigner > /home/postgres/database/notarysigner.bak'
    docker exec harbor-db sh -c 'pg_dump -U postgres notaryserver > /home/postgres/database/notaryserver.bak'
    docker cp harbor-db:/home/postgres/database harbor/
    #chmod 755 harbor/database
  }
  
  save_database() {
    cp -rf data/database harbor/
    #chmod 755 harbor/database
  }
  
  backup_database() {
    if [ "$kind" == "physical" ]; then
      save_database
    else
      export_database
    fi
    if [ "$pause" == "true" ]; then
      read -p "backup paused, continue?" tmpvar
    fi
  }
  
  backup_registry() {
    cp -rf data/registry harbor/
  }
  
  backup_chart_museum() {
    if [ -d data/chart_storage ]; then
      cp -rf data/chart_storage harbor/
    fi
  }
  
  backup_redis() {
    if [ -d data/redis ]; then
      cp -rf data/redis harbor/
    fi
  }
  
  backup_secret() {
    if [ -f data/secretkey ]; then
      cp data/secretkey harbor/secret/
    fi
    if [ -f data/defaultalias ]; then
      cp data/defaultalias harbor/secret/
    fi
    # location changed after 1.8.0
    if [ -d data/secret/keys/ ]; then
      cp -r data/secret/keys/ harbor/secret/
    fi
  }
  
  function backup() {
    create_dir
    backup_database
    backup_redis
    backup_registry
    backup_chart_museum
    backup_secret
    create_tarball $1
    clean_harbor
  }
  
  ##### restore
  extract_backup(){
    rm -rf data
    mkdir -p data/secret
    if [ -n "$1" ]; then
      tar xvf $1
    else
      tar xvf harbor.tgz
    fi
  }
  
  wait_for_db_ready() {
    chown -R 999:999 data/database
  
    TIMEOUT=12
    while [ $TIMEOUT -gt 0 ]; do
      docker exec harbor-db pg_isready | grep "accepting connections"
      if [ $? -eq 0 ]; then
        break
      fi
      TIMEOUT=$((TIMEOUT - 1))
      sleep 5
    done
    if [ $TIMEOUT -eq 0 ]; then
      echo "Harbor database cannot reach within one minute"
      exit 1
    fi
  }
  
  clean_database(){
    docker exec harbor-db psql -U postgres -d template1 -c "drop database registry;" 
    docker exec harbor-db psql -U postgres -d template1 -c "drop database postgres;"
    docker exec harbor-db psql -U postgres -d template1 -c "drop database notarysigner; "
    docker exec harbor-db psql -U postgres -d template1 -c "drop database notaryserver;"
  
    docker exec harbor-db psql -U postgres -d template1 -c "create database registry;"
    docker exec harbor-db psql -U postgres -d template1 -c "create database postgres;"
    docker exec harbor-db psql -U postgres -d template1 -c "create database notarysigner;"
    docker exec harbor-db psql -U postgres -d template1 -c "create database notaryserver;"
  }
  
  import_database() {
    docker exec harbor-db sh -c 'rm -rf /home/postgres/database'
    docker cp harbor/database harbor-db:/home/postgres/
    docker exec harbor-db sh -c 'psql -U postgres registry < /home/postgres/database/registry.bak'
    docker exec harbor-db sh -c 'psql -U postgres postgres < /home/postgres/database/postgres.bak'
    docker exec harbor-db sh -c 'psql -U postgres notarysigner < /home/postgres/database/notarysigner.bak'
    docker exec harbor-db sh -c 'psql -U postgres notaryserver < /home/postgres/database/notaryserver.bak'
  }
  
  load_database() {
    cp -r harbor/database/ data/
    chown -R 999:999 data/database
  }
  
  restore_database() {
    if [ -d ./harbor/database ]; then
      if [ "$kind" == "physical" ]; then
        load_database
        up
      else
        up
        wait_for_db_ready
        clean_database
        import_database
      fi
    fi
  }
  
  restore_registry() {
    if [ -d ./harbor/registry ]; then
      cp -r harbor/registry/ data/
      chown -R 10000:10000 data/registry
    fi
  }
  
  restore_redis() {
    if [ -d ./harbor/redis ]; then
      cp -r harbor/redis/ data/
      chown -R 999:999 data/redis
    fi
  }
  
  restore_chartmuseum() {
    if [ -d ./harbor/chart_storage ]; then
      cp -r harbor/chart_storage/ data/
      chown -R 10000:10000 data/chart_storage
    fi
  }
  
  restore_secret() {
    if [ -f harbor/secret/secretkey ]; then
      cp -f harbor/secret/secretkey data/secretkey 
    fi
    if [ -f harbor/secret/defaultalias ]; then
      cp -f harbor/secret/defaultalias data/secretkey 
    fi
    if [ -d harbor/secret/keys ]; then
      cp -r harbor/secret/keys data/secret/
    fi
  }
  
  restore() {
    extract_backup $1
    restore_redis
    restore_registry
    restore_chartmuseum
    restore_secret
    restore_database
    clean_harbor
  }
  
  ##### tools
  function push() {
    img="test.ksc.com/library/alpine"
    if [ -n "$1" ]; then
      img="$1"
    else
      docker tag alpine test.ksc.com/library/alpine
    fi
    docker push $img
  }
  
  function pull() {
    img="test.ksc.com/library/alpine"
    if [ -n "$1" ]; then
      img="$1"
    fi
    docker pull $img
  }
  
  function readonly() {
    ro="true"
    if [ -n "$1" ]; then
      ro="$1"
    fi
    curl -XPUT -u admin:Harbor12345 http://test.ksc.com/api/configurations -H "Content-Type: application/json" -d "{\"read_only\": $ro}"
  }
  
  function query() {
    id=1
    if [ -n "$1" ]; then
      id="$1"
    fi
    curl http://test.ksc.com/api/repositories?project_id=$id
  }
  
  function usage() {
    echo "===== USAGE ====="
    echo "./hbr [options] [action]"
    echo "==== OPTIONS ===="
    echo "-p: physical backup/restore for metadata"
    echo "-i: pause after harbor database backup finished"
    echo "-h: print help"
    echo "==== ACTION ====="
    echo "run: run a new harbor cluster"
    echo "up: make the harbor cluster up"
    echo "down: make the harbor cluster down"
    echo "clean: clean the harbor cluster"
    echo "backup: backup harbor data"
    echo "restore: restore harbor from backup"
    echo "push: push image to harbor"
    echo "pull: pull image from harbor"
    echo "readonly: set harbor to readonly mode"
    echo "query: query repositories using harbor api"
  }
  
  while getopts 'pih' OPT; do
    case $OPT in
    p)
      kind="physical"
      ;;
    i)
      pause="true"
      ;;
    h)
      usage
      exit 0
      ;;
    ?)
      usage
      exit 1
      ;;
    esac
  done
  
  shift $(($OPTIND - 1))
  
  case $1 in
  "run") run ;;
  "up") up ;;
  "down") down ;;
  "clean") clean ;;
  "backup") backup $2;;
  "restore") restore $2;;
  "push") push $2;;
  "pull") pull $2;;
  "query") query $2;;
  "readonly") readonly $2;;
  *) usage && exit 1 ;;
  esac
  ```


> docker-compose.yml中各组件的配置文件单独提供

## 可行性验证

验证备份恢复方案是否能够通过离线备在新的harbor集群中恢复数据.

### 构建Harbor集群

```bash
./hbr run
```

### 配置镜像仓环境

```bash
# 配置本机IP地址映射(不要用127.0.0.1)
echo "10.91.199.143 test.ksc.com" >> /etc/hosts   

# 配置/etc/docker/daemon.json
{
    "insecure-registries": ["test.ksc.com"],
    "registry-mirrors": [
        "http://test.ksc.com"
    ]
}

# docker登录harbor镜像仓
docker login test.ksc.com    # admin/Harbor12345
```

### 上传镜像到Harbor

```bash
# 前提: 本机已下载alpine:latest镜像
./hbr push
# 查看harbor中的镜像
./hbr query  # 也可以通过harbor webui查看镜像
```

### 备份Harbor数据

#### 全物理备份

```bash
./hbr -p backup
```

#### 半物理备份

```bash
./hbr backup
```

### 停止Harbor集群

```bash
./hbr down
```

### 恢复Harbor数据

#### 全物理恢复

```bash
./hbr -p restore <backup_filename>
```

#### 半物理恢复

```bash
./hbr restore <backup_filename>
```

### 启动Harbor集群

```bash
./hbr up
```

### 下载镜像到本地

```bash
./hbr pull
# 查看harbor中的镜像
./hbr query  # 也可以通过harbor webui查看镜像
```

### 设置Harbor为只读模式

```bash
./hbr readonly <true/false>   # 通过推送镜像验证, 或者打开harbor webui 配置管理查看
```

### 清除Harbor集群

```bash
./hbr clean
```

## 一致性验证

验证备份的数据库元数据和镜像仓/chart仓数据不一致时Harbor集群恢复后的状态.

可能导致不一致的原因是备份元数据和备份镜像/chart数据操作之间出现了数据增删操作.

以镜像仓数据为例, 需要验证两种情况: **有元数据但无镜像**和**有镜像但无元数据**.

#### 有元数据但无镜像

```bash
./hbr run
./hbr push
docker tag busybox test.ksc.com/library/busybox
./hbr push test.ksc.com/library/busybox
./hbr -i backup
# 备份暂停(数据库元数据备份已完成), 在webui上删除busybox镜像, 备份继续
./hbr down
./hbr restore <backup_filename>
./hbr pull test.ksc.com/library/busybox   # 拉取失败
./hbr clean    # 清理环境
```

*此情况下Harbor集群仍能正常启动, webui上仍可在相应项目下显示镜像名, 但tag数不存在. 拉取该镜像时, 会出现`Error response from daemon: manifest for test.ksc.com/library/busybox:latest not found: manifest unknown: manifest unknown`错误. 其他镜像不受影响.*

#### 有镜像但无元数据

```bash
./hbr run
./hbr push
./hbr -i backup
# 备份暂停(数据库元数据备份已完成), 再上传一个busybox镜像, 备份继续
docker tag busybox test.ksc.com/library/busybox
./hbr push test.ksc.com/library/busybox
./hbr down
./hbr restore <backup_filename>
./hbr pull test.ksc.com/library/busybox   # 拉取成功
./hbr clean     # 清理环境
```

*此情况下Harbor集群仍能正常启动, webui上相应项目下已不再显示镜像名, 但该镜像仍可正常拉取.*

**总结**: 当备份过程中元数据和镜像/charts数据的备份存在时间差时, 可能出现数据不一致的问题, 但不影响集群从备份中恢复. 为了避免数据不一致问题的出现, 应当在进行Harbor备份之前, 将集群状态设置为**只读模式**, 备份完成后再恢复状态.

> 只读模式可以在webui上设置, 也可以通过调用api设置
>
> 只读模式下, 无法上传, 删除或修改镜像, 但仍可提供镜像拉取服务

## 总结

1. 半物理和全物理方案都是可行的离线备份恢复方案, 高可用环境下需使用半物理方案
2. 为避免数据不一致问题, 应当将Harbor设置为只读模式后再进行数据备份
3. 本文档使用的编排和脚本文件仅作方案测试之用, 实际备份恢复脚本需视实际Harbor集群架构而定