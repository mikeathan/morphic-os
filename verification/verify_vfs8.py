import asyncio
from playwright.async_api import async_playwright

async def main():
    async with async_playwright() as p:
        browser = await p.chromium.launch()
        page = await browser.new_page(viewport={'width': 1280, 'height': 800})
        await page.goto('http://localhost:3000')
        # Wait for the map and vfs explorer to render
        await page.wait_for_selector('h2:has-text("Morphic Map")')
        await page.wait_for_selector('h2:has-text("VFS Explorer")')
        # Wait a moment for layout to settle
        await page.wait_for_timeout(2000)
        await page.screenshot(path='/app/vfs_verification8.png', full_page=True)
        await browser.close()

if __name__ == '__main__':
    asyncio.run(main())
