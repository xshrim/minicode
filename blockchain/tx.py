import core

w3 = core.myProvider("ws://127.0.0.1:8546", "test")
contractInfo = core.deploy(core.contractSourceCode, w3)
# contractInfo = core.getContract('./cinfo.json')
core.benchmark(w3, contractInfo, 2000, 'thread')