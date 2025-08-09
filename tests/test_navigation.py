import unittest
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager

class WebsiteNavigationTest(unittest.TestCase):
    def setUp(self):
        # Initialize Chrome driver
        self.driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()))
        # Maximize window to ensure dropdown is visible
        self.driver.maximize_window()
        # Load the website (adjust the path based on how you're serving the website)
        self.driver.get('file:///path/to/your/index.html')  # You'll need to update this path
        
    def tearDown(self):
        # Close the browser after each test
        self.driver.quit()
        
    def test_dropdown_visibility(self):
        """Test that dropdown menu items are initially hidden and appear on hover"""
        # Find the dropdown toggle element
        dropdown = self.driver.find_element(By.CLASS_NAME, "dropdown")
        dropdown_menu = self.driver.find_element(By.CLASS_NAME, "dropdown-menu")
        
        # Check that dropdown menu is initially not visible
        self.assertFalse(dropdown_menu.is_displayed(), "Dropdown menu should be hidden initially")
        
        # Hover over the dropdown
        webdriver.ActionChains(self.driver).move_to_element(dropdown).perform()
        
        # Wait for dropdown to become visible
        WebDriverWait(self.driver, 10).until(
            EC.visibility_of_element_located((By.CLASS_NAME, "dropdown-menu"))
        )
        
        # Verify dropdown is now visible
        self.assertTrue(dropdown_menu.is_displayed(), "Dropdown menu should be visible after hover")
        
    def test_dropdown_items(self):
        """Test that all dropdown items are present and have correct text"""
        # Expected dropdown items
        expected_items = ['Projects', 'Contact', 'Resources']
        
        # Hover over dropdown to make it visible
        dropdown = self.driver.find_element(By.CLASS_NAME, "dropdown")
        webdriver.ActionChains(self.driver).move_to_element(dropdown).perform()
        
        # Wait for dropdown menu to be visible
        WebDriverWait(self.driver, 10).until(
            EC.visibility_of_element_located((By.CLASS_NAME, "dropdown-menu"))
        )
        
        # Get all dropdown items
        dropdown_items = self.driver.find_elements(By.CSS_SELECTOR, ".dropdown-menu li a")
        
        # Verify number of items
        self.assertEqual(len(dropdown_items), len(expected_items), 
                        f"Expected {len(expected_items)} dropdown items, but found {len(dropdown_items)}")
        
        # Verify text of each item
        for item, expected_text in zip(dropdown_items, expected_items):
            self.assertEqual(item.text, expected_text, 
                           f"Expected dropdown item text '{expected_text}', but found '{item.text}'")
            
    def test_dropdown_links(self):
        """Test that dropdown links have correct href attributes"""
        # Expected href values
        expected_hrefs = ['#projects', '#contact', '#resources']
        
        # Hover over dropdown to make it visible
        dropdown = self.driver.find_element(By.CLASS_NAME, "dropdown")
        webdriver.ActionChains(self.driver).move_to_element(dropdown).perform()
        
        # Wait for dropdown menu to be visible
        WebDriverWait(self.driver, 10).until(
            EC.visibility_of_element_located((By.CLASS_NAME, "dropdown-menu"))
        )
        
        # Get all dropdown links
        dropdown_links = self.driver.find_elements(By.CSS_SELECTOR, ".dropdown-menu li a")
        
        # Verify href attributes
        for link, expected_href in zip(dropdown_links, expected_hrefs):
            actual_href = link.get_attribute('href')
            self.assertTrue(actual_href.endswith(expected_href), 
                          f"Expected href to end with '{expected_href}', but found '{actual_href}'")
            
    def test_dropdown_hover_styles(self):
        """Test that dropdown items change style on hover"""
        # Hover over dropdown to make it visible
        dropdown = self.driver.find_element(By.CLASS_NAME, "dropdown")
        webdriver.ActionChains(self.driver).move_to_element(dropdown).perform()
        
        # Wait for dropdown menu to be visible
        WebDriverWait(self.driver, 10).until(
            EC.visibility_of_element_located((By.CLASS_NAME, "dropdown-menu"))
        )
        
        # Get first dropdown item
        first_item = self.driver.find_element(By.CSS_SELECTOR, ".dropdown-menu li:first-child a")
        
        # Get initial background color
        initial_bg_color = first_item.value_of_css_property('background-color')
        
        # Hover over the item
        webdriver.ActionChains(self.driver).move_to_element(first_item).perform()
        
        # Wait for hover state
        import time
        time.sleep(0.5)  # Small delay to ensure hover styles are applied
        
        # Get background color after hover
        hover_bg_color = first_item.value_of_css_property('background-color')
        
        # Verify background color changed
        self.assertNotEqual(initial_bg_color, hover_bg_color, 
                          "Background color should change on hover")

if __name__ == '__main__':
    unittest.main() 