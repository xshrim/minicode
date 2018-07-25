import uuid, time, string, random, hashlib


class Token:
    idhash = "0x00"
    parenthash = "0x00"
    origincode = "0x00"
    statuscode = "0x00"
    typecode = '0x00'
    itemcode = "0x00"
    unitcode = "0x00"
    chancode = "0x00"
    keycode = "0x00"
    gaddr = "0x00"
    saddr = "0x00"
    caddr = "0x00"
    paddr = "0x00"
    snetwork = 0
    cnetwork = 0
    amount = 0
    gtimestamp = 0
    ftimestamp = 0
    etimestamp = 0
    ctimestamp = 0
    stimestamp = 0
    status = 0
    reserve = 0

    def __init__(self, idhash, parenthash, origincode, statuscode, typecode, itemcode, unitcode, chancode, keycode, gaddr, saddr,
                 caddr, paddr, snetwork, cnetwork, amount, gtimestamp, ftimestamp, etimestamp, ctimestamp, stimestamp, status,
                 reserve):
        self.idhash = idhash
        self.parenthash = parenthash
        self.origincode = origincode
        self.statuscode = statuscode
        self.typecode = typecode
        self.itemcode = itemcode
        self.unitcode = unitcode
        self.chancode = chancode
        self.keycode = keycode
        self.gaddr = gaddr
        self.saddr = saddr
        self.caddr = caddr
        self.paddr = paddr
        self.snetwork = snetwork
        self.cnetwork = cnetwork
        self.amount = amount
        self.gtimestamp = gtimestamp
        self.ftimestamp = ftimestamp
        self.etimestamp = etimestamp
        self.ctimestamp = ctimestamp
        self.stimestamp = stimestamp
        self.status = status
        self.reserve = reserve

    def __eq__(self, other):
        return self.idhash == other.idhash and self.parenthash == other.parenthash and self.origincode == other.origincode and self.statuscode == other.statuscode and self.typecode == other.typecode and self.itemcode == other.itemcode and self.unitcode == other.unitcode and self.chancode == other.chancode and self.keycode == other.keycode and self.gaddr == other.gaddr and self.saddr == other.saddr and self.caddr == other.caddr and self.paddr == other.padr and self.snetwork == other.cnetwork and self.amount == other.amount and self.gtimestamp == other.gtimestamp and self.ftimestamp == other.ftimestamp and self.etimestamp == other.etimestamp and self.ctimestamp == other.ctimestamp and self.stimestamp == other.stimestamp and self.status == other.status and self.reserve == other.reserve

    def __str__(self):
        return str({
            "idhash": self.idhash,
            "parenthash": self.parenthash,
            "origincode": self.origincode,
            "statuscode": self.statuscode,
            "typecode": self.typecode,
            "itemcode": self.itemcode,
            "unitcode": self.unitcode,
            "chancode": self.chancode,
            "keycode": self.keycode,
            "gaddr": self.gaddr,
            "saddr": self.saddr,
            "caddr": self.caddr,
            "paddr": self.paddr,
            "snetwork": self.snetwork,
            "cnetwork": self.cnetwork,
            "amount": self.amount,
            "gtimestamp": self.gtimestamp,
            "ftimestamp": self.ftimestamp,
            "etimestamp": self.etimestamp,
            "ctimestamp": self.ctimestamp,
            "stimestamp": self.stimestamp,
            "status": self.status,
            "reserve": self.reserve
        })


def randata(data, mode="default", num=10):
    if mode == "str":
        data.append(''.join(random.sample(string.ascii_letters + string.digits, num)))
    elif mode == "addr":
        data.append("0x" + str(hashlib.sha1(uuid.uuid1().bytes).hexdigest()))
    elif mode == "hash":
        # data.append("0x" + str(hashlib.sha256(uuid.uuid1().bytes).hexdigest()))
        data.append("0x" + str(uuid.uuid1().hex))
    elif mode == "small":
        data.append(random.randint(100, 999))
    elif mode == "large":
        data.append(random.randint(0, 9999999999999999))
    elif mode == "timestamp":
        data.append(int(time.time()) + random.randint(-1000000, 1000000))
    else:
        random.shuffle(data)
    return random.choice(data)
    # return random.sample(data, 1)[0]


def newToken():
    try:
        metadata = [
            "0x00" + "0" * 30, "0x01" + "0" * 30, "0x10" + "0" * 30, "0x11" + "0" * 30, "credit", "tender", "yuan", "ge", "jian",
            "xiang"
        ]

        idhash = randata([], "hash")
        parenthash = "0x00" + "0" * 30
        origincode = randata(metadata[:4], "hash")
        statuscode = randata(metadata[:4], "hash")
        # typecode = str(hashlib.sha256(randata(metadata[4:6], "str", 6).encode('utf8')).hexdigest())
        typecode = randata(metadata[4:6], "str", 6).encode('utf8')
        itemcode = randata([], "hash")
        # unitcode = str(hashlib.sha256(randata(metadata[6:], "str", 6).encode('utf8')).hexdigest())
        unitcode = randata(metadata[6:], "str", 6).encode('utf8')
        chancode = randata([], "str", 6).encode('utf8')
        keycode = "0x00" + "0" * 30
        gaddr = randata([], "addr")
        saddr = randata([], "addr")
        caddr = randata([saddr], "addr")
        paddr = randata([], "addr")
        snetwork = randata([], "small")
        cnetwork = snetwork
        amount = randata([], "large")
        gtimestamp = randata([], "timestamp")
        ftimestamp = randata([gtimestamp], "timestamp")
        etimestamp = randata([gtimestamp], "timestamp")
        ctimestamp = randata([gtimestamp], "timestamp")
        stimestamp = randata([gtimestamp], "timestamp")
        status = randata([0, 1, 2, 6], "none")
        reserve = 0

        return Token(
            str(idhash), str(parenthash), str(origincode), str(statuscode), str(typecode), str(itemcode), str(unitcode),
            str(chancode), str(keycode), str(gaddr), str(saddr), str(caddr), str(paddr), int(snetwork), int(cnetwork),
            int(amount), int(gtimestamp), int(ftimestamp), int(etimestamp), int(ctimestamp), int(stimestamp), int(status),
            int(reserve))
    except Exception as ex:
        print(str(ex))
        return None


print("aa".encode("utf8"))
