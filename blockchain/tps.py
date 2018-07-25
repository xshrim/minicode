import core
# from web3.auto import w3

w3 = core.myProvider("ws://127.0.0.1:8546")
# contractInfo = deploy(contractSourceCode, w3)
# contractInfo = core.getContract('./cinfo.json')
# core.invoke(w3, contractInfo, 1000, 'thread')
core.monitor(w3)
