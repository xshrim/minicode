#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""Download a webpage as a PDF."""

from selenium import webdriver


def download(driver, target_path):
    """Download the currently displayed page to target_path."""

    def execute(script, args):
        driver.execute('executePhantomScript', {'script': script, 'args': args})

    # hack while the python interface lags
    driver.command_executor._commands['executePhantomScript'] = ('POST', '/session/$sessionId/phantom/execute')
    # set page format
    # inside the execution script, webpage is "this"
    page_format = 'this.paperSize = {format: "A4", orientation: "portrait" };'
    execute(page_format, [])

    # render current page
    render = '''this.render("{}")'''.format(target_path)
    execute(render, [])


if __name__ == '__main__':
    driver = webdriver.PhantomJS('phantomjs')   # phantomjs已过时, 推荐chromedriver
    driver.get('https://blog.csdn.net/wangshubo1989/article/details/55102275')
    download(driver, "save_me.pdf")
