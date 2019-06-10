import re, sys
from pyh import *


class Document:

    head = ""
    name = ""
    auth = ""
    date = ""
    proj = ""
    desc = ""

    def __init__(self, head, name, auth, date, proj, desc):
        self.head = head
        self.name = name
        self.auth = auth
        self.date = date
        self.proj = proj
        self.desc = desc

    def __str__(self):
        return str({
            "head": self.head,
            "name": self.name,
            "auth": self.auth,
            "date": self.date,
            "proj": self.proj,
            "desc": self.desc
        })


class Contract:

    head = ""
    desc = ""

    kind = ""
    name = ""

    structs = []
    events = []
    functions = []

    def __init__(self, head, desc, kind, name, structs, events, functions):
        self.head = head
        self.desc = desc
        self.kind = kind
        self.name = name
        self.structs = structs
        self.events = events
        self.functions = functions

    def __str__(self):
        return str({
            "head": self.head,
            "desc": self.desc,
            "kind": self.kind,
            "name": self.name,
            "structs": self.structs,
            "events": self.events,
            "functions": self.functions
        })


class Struct:

    head = ""
    desc = ""

    kind = ""
    name = ""

    mems = []

    def __init__(self, head, desc, kind, name, mems):
        self.head = head
        self.desc = desc
        self.kind = kind
        self.name = name
        self.mems = mems

    def __str__(self):
        return str({"head": self.head, "desc": self.desc, "kind": self.kind, "name": self.name, "mems": self.mems})


class Function:
    head = ""
    desc = ""

    visb = ""
    sync = ""
    epay = ""
    modf = ""
    kind = ""
    name = ""

    args = []
    rets = []

    def __init__(self, head, desc, visb, sync, epay, modf, kind, name, args, rets):
        self.head = head
        self.desc = desc
        self.visb = visb
        self.sync = sync
        self.epay = epay
        self.modf = modf
        self.kind = kind
        self.name = name
        self.args = args
        self.rets = rets

    def __str__(self):
        return str({
            "head": self.head,
            "desc": self.desc,
            "visb": self.visb,
            "sync": self.sync,
            "epay": self.epay,
            "modf": self.modf,
            "kind": self.kind,
            "name": self.name,
            "args": self.args,
            "rets": self.rets
        })


class Event:
    head = ""
    desc = ""

    kind = ""
    name = ""

    args = []

    def __init__(self, head, desc, kind, name, args):
        self.head = head
        self.desc = desc
        self.kind = kind
        self.name = name
        self.args = args

    def __str__(self):
        return str({"head": self.head, "desc": self.desc, "kind": self.kind, "name": self.name, "args": self.args})


def parse(fpath):
    document = []
    with open(fpath, "r") as rf:
        filehead = "文件描述"
        filename = ""
        fileauth = ""
        filedate = ""
        fileproj = ""
        filedesc = ""

        contracts = []

        source = rf.read()

        rfilename = re.findall(r'//@filename\s*:\s*(.*)', source)
        if rfilename is not None and len(rfilename) > 0:
            filename = str(rfilename[0]).strip()
        rfileauth = re.findall(r'//@author\s*:\s*(.*)', source)
        if rfileauth is not None and len(rfileauth) > 0:
            fileauth = str(rfileauth[0]).strip()
        rfiledate = re.findall(r'//@date\s*:\s*(.*)', source)
        if rfiledate is not None and len(rfiledate) > 0:
            filedate = str(rfiledate[0]).strip()
        rfileproj = re.findall(r'//@project\s*:\s*(.*)', source)
        if rfileproj is not None and len(rfileproj) > 0:
            fileproj = str(rfileproj[0]).strip()
        rfiledesc = re.findall(r'//@describe\s*:\s*(.*)', source)
        if rfiledesc is not None and len(rfiledesc) > 0:
            filedesc = str(rfiledesc[0]).strip()
        document.append(Document(filehead, filename, fileauth, filedate, fileproj, filedesc))

        parts = re.split(r'contract', source)
        if len(parts) < 2:
            return contracts
        for i in range(1, len(parts)):
            content = re.findall(r'(/\*@.*?\*/)', parts[i - 1], re.DOTALL)[-1] + "\ncontract " + parts[i]
            # print(content)
            # break

            # sys.exit()
            # for contract in re.findall(r'(/\*@.*?\*/.*?contract\s+\w+\s*?\{.*\})', source, re.DOTALL):
            # content = contract
            chead = ""
            cdesc = ""

            ckind = ""
            cname = ""

            cstructs = []
            cfunctions = []
            cevents = []

            rs = re.search(r'(/\*@.*?\*/)', content, re.DOTALL)
            if rs is not None:
                comment = rs.group(1)
                for rs in re.findall(r'/\*@(.*)', comment):
                    chead = rs.strip()
                for rs in re.findall(r'//@desc(.*)', comment):
                    cdesc += str(rs).strip()
            rs = re.search(r'(contract)\s+(\w+)\s*?(\{.*\})', content, re.DOTALL)
            if rs is not None:
                ckind = rs.group(1).strip()
                cname = rs.group(2).strip()
                cbody = rs.group(3).strip()
                print("%" * 40 + (str(ckind) + " <--> " + str(cname)).center(20, " ") + "%" * 40)
                if cbody is not None:
                    content = cbody
                    for body in re.findall(r'(/\*@.*?\*/.*?(struct|function|event)\s+.*?(\{.*?\}|;))', content, re.DOTALL):
                        content = body[0]
                        if re.search(r'(/\*@.*?\*/.*?struct\s+.*?\{.*?\})', content, re.DOTALL) is not None:
                            shead = ""
                            sdesc = ""

                            skind = ""
                            sname = ""

                            smems = []

                            rs = re.search(r'(/\*@.*?\*/)', content, re.DOTALL)
                            if rs is not None:
                                comment = rs.group(1)
                                # print(comment)
                                for rs in re.findall(r'/\*@(.*)', comment):
                                    shead = rs.strip()
                                for rs in re.findall(r'//@desc(.*)', comment):
                                    sdesc += str(rs).strip()
                            rs = re.search(r'(struct)\s+(.*)\s*\{', content)
                            if rs is not None:
                                skind = rs.group(1).strip()
                                sname = rs.group(2).strip()
                            print("#" * 20 + (str(skind) + " <--> " + str(sname)).center(40, " ") + "#" * 20)

                            rs = re.search(r'\{(.*)\}', content, re.DOTALL)
                            if rs is not None:
                                for rs in re.findall(r'(.*);(.*)', rs.group(1)):
                                    svar = str(rs[0]).strip()
                                    svcm = str(rs[1]).strip()
                                    if svar[0] == '/':
                                        continue
                                    svartype, svarname = svar.split(' ')
                                    svartype = svartype.strip()
                                    svarname = svarname.strip()
                                    if "address" in svartype:
                                        svarlens = "20"
                                    elif "int" in svartype:
                                        svardig = svartype[svartype.index('int') + 3:]
                                        if svardig == "":
                                            svarlens = "32"
                                        else:
                                            if '[' in svardig:
                                                svarlens = str(int(int(svardig[:svardig.index('[')]) / 8))
                                            else:
                                                svarlens = str(int(int(svardig) / 8))
                                    elif "byte" in svartype:
                                        if "bytes" in svartype:
                                            svardig = svartype[svartype.index('bytes') + 5:]
                                            if svardig == "":
                                                svarlens = "0"
                                            else:
                                                if '[' in svardig:
                                                    svarlens = str(svardig[:svardig.index('[')])
                                                svarlens = str(svardig.strip())
                                        else:
                                            svarlens = "1"
                                    else:
                                        svarlens = "0"
                                    svcm = svcm.replace('/', '').strip()
                                    if '(' in svcm:
                                        svcmdesc = svcm[:svcm.index('(')]
                                        svcmcomm = svcm[svcm.index('('):]
                                    else:
                                        svcmdesc = svcm
                                        svcmcomm = ""
                                    smems.append([svarname, svartype, svarlens, svcmdesc, svcmcomm])
                            cstructs.append(Struct(shead, sdesc, skind, sname, smems))

                        elif re.search(r'(/\*@.*?\*/.*?event\s+.*?;)', content, re.DOTALL) is not None:
                            ehead = ""
                            edesc = ""
                            cargs = {}

                            ekind = ""
                            ename = ""

                            iargs = []

                            rs = re.search(r'(/\*@.*?\*/)', content, re.DOTALL)
                            if rs is not None:
                                comment = rs.group(1)
                                for rs in re.findall(r'/\*@(.*)', comment):
                                    ehead = rs.strip()
                                for rs in re.findall(r'//@desc(.*)', comment):
                                    edesc += str(rs).strip()
                                for rs in re.findall(r'//@param(.*)', comment):
                                    arg, desc = rs.split(':')
                                    cargs[arg.strip()] = desc.strip()

                            rs = re.search(r'(event)\s+(\w+)\s*\((.*?)\).*?;', content)
                            if rs is not None:
                                # ekind, ename, eargs = re.search(r'(event)\s+(\w+)\s*\((.*?)\).*?;', content).groups()
                                ekind = rs.group(1).strip()
                                ename = rs.group(2).strip()
                                eargs = rs.group(3).strip()
                                print("+" * 20 + (str(ekind) + " <--> " + str(ename)).center(40, " ") + "+" * 20)

                                if eargs is not None and eargs != "":
                                    for arg in eargs.split(','):
                                        # print(arg)
                                        arg = arg.strip()
                                        if "indexed" in arg:
                                            argindx = "indexed"
                                            argtype, _, argname = arg.split(' ')
                                        else:
                                            argindx = ""
                                            argtype, argname = arg.split(' ')
                                        # argtype, argname = arg.split(' ')
                                        argtype = argtype.strip()
                                        argname = argname.strip()
                                        argdescs = str(cargs[argname]).strip()
                                        if '(' in argdescs:
                                            argdesc = argdescs[:argdescs.index('(')]
                                            argcomm = argdescs[argdescs.index('('):]
                                        else:
                                            argdesc = argdescs
                                            argcomm = ""
                                        if "address" in argtype:
                                            arglens = "20"
                                        elif "int" in argtype:
                                            argdig = argtype[argtype.index('int') + 3:]
                                            if argdig == "":
                                                arglens = "32"
                                            else:
                                                if '[' in argdig:
                                                    arglens = str(int(int(argdig[:argdig.index('[')]) / 8))
                                                else:
                                                    arglens = str(int(int(argdig) / 8))
                                        elif "byte" in argtype:
                                            if "bytes" in argtype:
                                                argdig = argtype[argtype.index('bytes') + 5:]
                                                if argdig == "":
                                                    arglens = "0"
                                                else:
                                                    if '[' in argdig:
                                                        arglens = str(argdig[:argdig.index('[')])
                                                    else:
                                                        arglens = str(argdig.strip())
                                            else:
                                                arglens = "1"
                                        else:
                                            arglens = "0"
                                        iargs.append([argname, argtype, arglens, argdesc, argcomm, argindx])

                            cevents.append(Event(ehead, edesc, ekind, ename, iargs))

                        elif re.search(r'(/\*@.*?\*/.*?function\s+.*?\{.*?\})', content, re.DOTALL) is not None:
                            fhead = ""
                            fdesc = ""
                            cargs = {}
                            crets = []

                            fvisb = ""
                            fsync = ""
                            fepay = ""
                            fmodf = ""
                            fkind = ""
                            fname = ""

                            iargs = []
                            irets = []

                            rs = re.search(r'(/\*@.*?\*/)', content, re.DOTALL)
                            if rs is not None:
                                comment = rs.group(1)
                                for rs in re.findall(r'/\*@(.*)', comment):
                                    fhead = rs.strip()
                                for rs in re.findall(r'//@desc(.*)', comment):
                                    fdesc += str(rs).strip()
                                for rs in re.findall(r'//@param(.*)', comment):
                                    arg, desc = rs.split(':')
                                    cargs[arg.strip()] = desc.strip()
                                for rs in re.findall(r'//@return(.*)', comment):
                                    crets.append(rs.split(':')[1].strip())
                            rs = re.search(r'(function\s+.*?\{)', content)
                            if rs is not None:
                                function = rs.group(1)
                                # print(function)

                                if " public " in function:
                                    fvisb = "public"
                                elif " private " in function:
                                    fvisb = "private"
                                elif " internal " in function:
                                    fvisb = "internal"
                                if " view " in function or " pure " in function:
                                    fsync = "call"
                                else:
                                    fsync = "send"
                                if " payable " in function:
                                    fepay = "payable"
                                else:
                                    fepay = ""

                                rs = re.search(r'(function)\s+(\w+)\s*\((.*?)\)(.*)', function)
                                if rs is not None:
                                    fkind = rs.group(1).strip()
                                    fname = rs.group(2).strip()
                                    fargs = rs.group(3).strip()
                                    frets = rs.group(4).strip()
                                    print("*" * 20 + (str(fkind) + " <--> " + str(fname)).center(40, " ") + "*" * 20)
                                    # fkind, fname, fargs, frets = re.search(r'(function)\s+(\w+)\s*\((.*?)\)(.*)', function).groups()
                                    if "return" in frets:
                                        fmodf = frets[:frets.index("return")]
                                        fmodf = fmodf.replace("public", "").replace("private",
                                                                                    "").replace("internal", "").replace(
                                                                                        "view", "").replace("pure", "").replace("payable", "").strip()
                                    else:
                                        fmodf = ""
                                    if "(" in frets and ")" in frets:
                                        frets = re.search(r'\((.*?)\)', frets).group(1).strip()
                                    else:
                                        frets = ""
                                    if fargs is not None and fargs != "":
                                        for arg in fargs.split(','):
                                            argtype, argname = arg.strip().split(' ')
                                            argtype = argtype.strip()
                                            argname = argname.strip()
                                            argdescs = str(cargs[argname]).strip()
                                            if '(' in argdescs:
                                                argdesc = argdescs[:argdescs.index('(')]
                                                argcomm = argdescs[argdescs.index('('):]
                                            else:
                                                argdesc = argdescs
                                                argcomm = ""
                                            if "address" in argtype:
                                                arglens = "20"
                                            elif "int" in argtype:
                                                argdig = argtype[argtype.index('int') + 3:]
                                                if argdig == "":
                                                    arglens = "32"
                                                else:
                                                    if '[' in argdig:
                                                        arglens = str(int(int(argdig[:argdig.index('[')]) / 8))
                                                    else:
                                                        arglens = str(int(int(argdig) / 8))
                                            elif "byte" in argtype:
                                                if "bytes" in argtype:
                                                    argdig = argtype[argtype.index('bytes') + 5:]
                                                    if argdig == "":
                                                        arglens = "0"
                                                    else:
                                                        if '[' in argdig:
                                                            arglens = str(argdig[:argdig.index('[')])
                                                        else:
                                                            arglens = str(argdig.strip())
                                                else:
                                                    arglens = "1"
                                            else:
                                                arglens = "0"
                                            iargs.append([argname, argtype, arglens, argdesc, argcomm])

                                    if frets is not None and frets != "":
                                        rets = frets.split(',')
                                        for i in range(len(rets)):
                                            # print(ret)
                                            ret = rets[i].strip()
                                            if " " in ret:
                                                rettype, retname = ret.split(' ')
                                            else:
                                                rettype = ret
                                                retname = ""
                                            rettype = rettype.strip()
                                            retname = retname.strip()
                                            retdescs = str(crets[i]).strip()
                                            if "(" in retdescs:
                                                retdesc = retdescs[:retdescs.index('(')]
                                                retcomm = retdescs[retdescs.index('('):]
                                            else:
                                                retdesc = retdescs
                                                retcomm = ""
                                            if "address" in rettype:
                                                retlens = "20"
                                            elif "int" in rettype:
                                                retdig = rettype[rettype.index("int") + 3:]
                                                if retdig == "":
                                                    retlens = "32"
                                                else:
                                                    if "[" in retdig:
                                                        retlens = str(int(int(retdig[:retdig.index("[")]) / 8))
                                                    else:
                                                        retlens = str(int(int(retdig) / 8))
                                            elif "byte" in rettype:
                                                if "bytes" in rettype:
                                                    retdig = rettype[rettype.index("bytes") + 5:]
                                                    if retdig == "":
                                                        retlens = "0"
                                                    else:
                                                        if "[" in retdig:
                                                            retlens = str(retdig[:retdig.index("[")])
                                                        retlens = str(retdig.strip())
                                                else:
                                                    retlens = "1"
                                            else:
                                                retlens = "0"
                                            irets.append([retname, rettype, retlens, retdesc, retcomm])

                            cfunctions.append(Function(fhead, fdesc, fvisb, fsync, fepay, fmodf, fkind, fname, iargs, irets))
                        # print("#" * 50)
            contracts.append(Contract(chead, cdesc, ckind, cname, cstructs, cevents, cfunctions))
            document.append(contracts)
            # print("%" * 50)
    return document


'''
for con in parse("./RS.sol"):
    print(len(con.structs), len(con.events), len(con.functions))
    print("=" * 100)
    for struct in con.structs:
        print(struct)
    print("=" * 100)
    for event in con.events:
        print(event)
    print("=" * 100)
    for function in con.functions:
        print(function)
    print("=" * 100)
sys.exit()
'''


def pretty(document):
    page = PyH('API')
    page.addCSS('mystyle.css', 'yourstyle.css')
    page.addJS('myjs.js', 'yourjs.js')
    title = page << h1('接口文档', cl='titlehead')

    abstractDiv = page << div(cl='abstractdiv')
    abstractDiv << h2('摘要', cl='absthead')
    abstractDiv << div(cl='abstdesc')

    catalogDiv = page << div(cl='catalogdiv')
    catalogDiv << h2('目录', cl='catahead')
    catalogUL = catalogDiv << ul(cl='cataul clear')
    catalogLI = catalogUL << li('目录')

    catastrUL = catalogLI << ul(cl='catastrul clear')
    cataeveUL = catalogLI << ul(cl='cataeveul clear')
    catafunUL = catalogLI << ul(cl='catafunul clear')

    catastrLI = catastrUL << li('结构')
    cataeveLI = cataeveUL << li('事件')
    catafunLI = catafunUL << li('函数')

    cataistrUL = catastrLI << ul(cl='cataistrul clear')

    cataieveUL = cataeveLI << ul(cl='cataieveul clear')

    cataifunUL = catafunLI << ul(cl='cataifunul clear')

    head = ""
    name = ""
    auth = ""
    date = ""
    proj = ""
    desc = ""

    dochead = document[0]
    '''
    abstractDiv << h5('文件')
    abstractDiv << p(dochead.name)
    abstractDiv << h5('作者')
    abstractDiv << p(dochead.auth)
    abstractDiv << h5('日期')
    abstractDiv << p(dochead.date)
    abstractDiv << h5('项目')
    abstractDiv << p(dochead.proj)
    abstractDiv << h5('描述')
    abstractDiv << p(dochead.desc)
    '''

    abstractDiv << h5('文件' + ': ' + dochead.name)
    abstractDiv << h5('作者' + ': ' + dochead.auth)
    abstractDiv << h5('日期' + ': ' + dochead.date)
    abstractDiv << h5('项目' + ': ' + dochead.proj)
    abstractDiv << h5('描述' + ': ' + dochead.desc)

    doccont = document[1]
    for doc in doccont:
        contractDiv = page << div(cl='contractdiv')
        contractDiv << h2('合约 - ' + doc.name)
        concommDiv = contractDiv << div(cl='concomm')
        concommDiv << div(cl='conheaddiv') << p(doc.head)
        concommDiv << div(cl='condescdiv') << p(doc.desc)

        comcontDiv = contractDiv << div(cl='concont')

        structDiv = comcontDiv << div(cl='structdiv')
        structDiv << h3('合约结构')

        eventDiv = comcontDiv << div(cl='eventdiv')
        eventDiv << h3('合约事件')

        functionDiv = comcontDiv << div(cl='functiondiv')
        functionDiv << h3('合约函数')

        for struct in doc.structs:
            structDiv << div(
                cl='strtitlediv', id='struct_' + struct.name) << h4('结构 - ' + struct.name) << a(
                    '[目录]', href='#' + 'cata_' + struct.name)
            cataistrUL << li(cl='catastrli') << a(
                href='#struct_' + struct.name, id='cata_' + struct.name,
                title=struct.desc) << span(struct.name + ' [' + struct.head + ']')
            strcommDiv = structDiv << div(cl='strcomm')
            strcommDiv << div(cl='strheaddiv') << h5(struct.head)
            strcommDiv << div(cl='strdescdiv') << p(struct.desc)

            strcontDiv = structDiv << div(cl='strcont')

            strTable = strcontDiv << table(cl='strtable', border="1")
            strTable << caption(cl='strtablecap') << span(struct.name)
            strTable << tr(cl='strheadtr') << th('变量索引') + th('变量名称') + th('变量类型') + th('变量长度') + th('变量描述') + th('变量说明')

            for idx, item in enumerate(struct.mems):

                strTable << tr(cl='strconttr') << td(str(idx)) + td(str(item[0])) + td(str(
                    item[1])) + td(str(item[2]) if str(item[2]) != "0" else "-") + td(str(item[3])) + td(str(item[4]))

            # print(struct)

        for event in doc.events:
            eventDiv << div(
                cl='evetitlediv', id='event_' + event.name) << h4('事件 - ' + event.name) << a(
                    '[目录]', href='#' + 'cata_' + event.name)
            cataieveUL << li(cl='cataeveli') << a(
                href='#event_' + event.name, id='cata_' + event.name,
                title=event.desc) << span(event.name + ' [' + event.head + ']')
            evecommDiv = eventDiv << div(cl='evecomm')
            evecommDiv << div(cl='eveheaddiv') << h5(event.head)
            evecommDiv << div(cl='evedescdiv') << p(event.desc)

            evecontDiv = eventDiv << div(cl='evecont')

            eveTable = evecontDiv << table(cl='evetable', border="1")
            eveTable << caption(cl='evetablecap') << span(event.name)
            eveTable << tr(cl='eveheadtr') << th('变量索引') + th('变量名称') + th('变量类型') + th('变量长度') + th('变量描述') + th('变量说明')

            for idx, item in enumerate(event.args):
                eveTable << tr(cl='eveconttr') << td(str(idx)) + td(str(item[0])) + td(str(
                    item[1])) + td(str(item[2]) if str(item[2]) != "0" else "-") + td(str(item[3])) + td(
                        (str(item[5]) + ' ' + str(item[4])).strip())
            # print(event)

        for function in doc.functions:
            # print(function)
            funcqual = function.visb + '#' + function.sync
            if function.epay != '':
                funcqual += '#' + function.epay
            if function.modf != '':
                funcqual += '#' + function.modf
            functionDiv << div(
                cl='funtitlediv', id='function_' + function.name) << h4('函数 - {' + funcqual + '} ' + function.name) << a(
                    '[目录]', href='#' + 'cata_' + function.name)
            cataifunUL << li(cl='catafunli ' + function.visb + ' ' + function.sync) << a(
                href='#function_' + function.name, id='cata_' + function.name,
                title=function.desc) << span(function.name + ' [' + function.head + ']')
            funcommDiv = functionDiv << div(cl='funcomm')
            funcommDiv << div(cl='funheaddiv') << h5(function.head)
            funcommDiv << div(cl='fundescdiv') << p(function.desc)

            funcontDiv = functionDiv << div(cl='funcont')

            funTable = funcontDiv << table(cl='funtable', border="1")
            funTable << caption(cl='funtablecap') << span(function.name)
            funTable << tr(cl='funheadtr') << th('变量索引') + th('变量名称') + th('变量类型') + th('变量长度') + th('变量描述') + th('变量说明')

            funTable << tr(cl='funargtr') << th('参数', colspan="6", align='left')
            for idx, item in enumerate(function.args):
                funTable << tr(cl='funconttr') << td(str(idx)) + td(str(item[0])) + td(str(
                    item[1])) + td(str(item[2]) if str(item[2]) != "0" else "-") + td(str(item[3])) + td(str(item[4]))

            funTable << tr(cl='funrettr') << th('返回值', colspan="6", align='left')
            for idx, item in enumerate(function.rets):
                funTable << tr(cl='funconttr') << td(str(idx)) + td(str(item[0])) + td(str(
                    item[1])) + td(str(item[2]) if str(item[2]) != "0" else "-") + td(str(item[3])) + td(str(item[4]))

    page.printOut('./api.html')


pretty(parse("./RS.sol"))
