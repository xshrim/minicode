import json
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import Select
from selenium.webdriver.common.action_chains import ActionChains

appState = {
    "recentDestinations": [{
        "id": "Save as PDF",
        "origin": "local"
    }],
    "selectedDestinationId": "Save as PDF",
    "version": 2
}

profile = {
    'printing.print_preview_sticky_settings.appState': json.dumps(appState),
    'savefile.default_directory': "/home/xshrim/work/"
}

chrome_options = webdriver.ChromeOptions()
#chrome_options.add_experimental_option('prefs', profile)
#chrome_options.add_argument('--kiosk-printing')
#chrome_options.add_argument('headless')
#driver.maximize_window()

driver = webdriver.Chrome(options=chrome_options)
# ??????????
driver.set_page_load_timeout(10)
try:
    driver.get('https://www.cnblogs.com/fnng/p/3183777.html')
except Exception as ex:
    print('??????time out after 10 seconds when loading page??????')
    # ????????????????js?stop?????????
    driver.execute_script("window.stop()")

element = driver.find_element_by_tag_name('body')
element.send_keys(Keys.CONTROL + 'a')
driver.execute_script('window.print();')

#ActionChains(driver).move_to_element(element).send_keys(Keys.CONTROL, "p").perform()
print("YES")
print(driver.find_elements_by_css_selector("#post")[0].text)
#.send_keys(Keys.CONTROL, "a")
print("OK")
driver.execute_script('window.print();')

select = Select(driver.find_element_by_id('test-make-box'))
select.select_by_value('Test')
driver.find_element_by_css_selector('.btn.btn-primary.btn--search').click()

#print(driver.title)
#driver.find_element_by_id('kw').send_keys('测试')

driver.quit()