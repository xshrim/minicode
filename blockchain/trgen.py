import uuid, time, string, random, hashlib


class Record:
    sn = ""
    txtype = ""
    callfunc = ""
    txhash = "0x" + "0" * 32
    ifhash = "0x" + "0" * 32
    channel = "0x" + "0" * 32
    tonetwork = 0
    timestamp = 0
    sender = "0x" + "0" * 40
    payer = "0x" + "0" * 40
    payee = "0x" + "0" * 40
    splitids = []
    payids = []
    backids = []
    infos = ""

    def __init__(self, sn, txtype, callfunc, txhash, ifhash, channel, tonetwork, timestamp, sender, payer, payee, splitids,
                 payids, backids, infos):
        self.sn = sn
        self.txtype = txtype
        self.callfunc = callfunc
        self.txhash = txhash
        self.ifhash = ifhash
        self.channel = channel
        self.tonetwork = tonetwork
        self.timestamp = timestamp
        self.sender = sender
        self.payer = payer
        self.payee = payee
        self.splitids = splitids
        self.payids = payids
        self.backids = backids
        self.infos = infos

    def __str__(self):
        return str({
            "sn": self.sn,
            "txtype": self.txtype,
            "callfunc": self.callfunc,
            "txhash": self.txhash,
            "ifhash": self.ifhash,
            "channel": self.channel,
            "tonetwork": self.tonetwork,
            "timestamp": self.timestamp,
            "sender": self.sender,
            "payer": self.payer,
            "payee": self.payee,
            "splitids": self.splitids,
            "payids": self.payids,
            "backids": self.backids,
            "infos": self.infos
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


def newRecord():
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
        typecode = randata(metadata[4:6], "str", 6)
        itemcode = randata([], "hash")
        # unitcode = str(hashlib.sha256(randata(metadata[6:], "str", 6).encode('utf8')).hexdigest())
        unitcode = randata(metadata[6:], "str", 6)
        saddr = randata([], "addr")
        caddr = randata([saddr], "addr")
        snetwork = randata([], "small")
        cnetwork = snetwork
        amount = randata([], "large")
        gtimestamp = randata([], "timestamp")
        ftimestamp = randata([gtimestamp], "timestamp")
        etimestamp = randata([gtimestamp], "timestamp")
        ctimestamp = randata([gtimestamp], "timestamp")
        stimestamp = randata([gtimestamp], "timestamp")
        status = randata([0, 1, 2, 6], "none")

        return Token(
            str(idhash), str(parenthash), str(origincode), str(statuscode), str(typecode), str(itemcode), str(unitcode),
            str(saddr), str(caddr), int(snetwork), int(cnetwork), int(amount), int(gtimestamp), int(ftimestamp), int(etimestamp),
            int(ctimestamp), int(stimestamp), int(status))
    except Exception as ex:
        print(str(ex))
        return None
