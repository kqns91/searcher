import requests
from bs4 import BeautifulSoup
import json
import re
import datetime
from itertools import count
import os
from opensearchpy import OpenSearch
from aiohttp import web
import logging

base = 'https://www.nogizaka46.com'
open_date = '2011.11.11 00:00'

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)
handler = logging.StreamHandler()
handler.setLevel(logging.INFO)
formatter = logging.Formatter('%(asctime)s %(levelname)s %(message)s')
handler.setFormatter(formatter)
logger.addHandler(handler)

# mm.dd HH:MM => YYYY.mm.dd HH:MM
def add_year(date_string):
    now = datetime.datetime.now()

    month, day = map(int, date_string[:5].split("."))

    if datetime.datetime(now.year, month, day) > now:
        year = now.year - 1
    else:
        year = now.year

    date_str = f"{year:04d}.{month:02d}.{day:02d} {date_string[-5:]}"

    return date_str

def fetch_page(page_url):
    try:
        res = requests.get(page_url)
        res.raise_for_status()
        soup = BeautifulSoup(res.text, "html.parser")
    except requests.exceptions.RequestException as e:
        raise ValueError(f'failed to request: {e}')
    return soup

def get_updated_member_urls(url):
  try:
    soup = fetch_page(url)
  except Exception as e:
    raise ValueError(f'failed to fetch_page: {e}')

  members = []

  for option_tag in soup.find_all('option'):
      blog_path = option_tag.get('value')

      if blog_path.startswith('/s/n46/diary/MEMBER/list'):
        pattern = r'(.+?)\((.+?) 更新\)'
        result = re.search(pattern, option_tag.text)

        if result:
          name = result.group(1)
          update_time = add_year(result.group(2))

          member = {
            'name': name,
            'blog_url': base + blog_path,
            'update_time': update_time
          }

          members.append(member)

  try:
    with open("previous.json", "r") as f:
      previous = json.loads(f.read())
  except Exception as e:
    raise ValueError(f'failed to loads json: {e}')

  updated_members_urls = []
  new_members = []

  for _member in members:
    _name = _member['name']
    _update_time = _member['update_time']
    _blog_url = _member['blog_url']

    for i, member in enumerate(previous):
      name = member['name']
      update_time = member['update_time']

      if _name == name:
        if _update_time == update_time:
          break

        blog_info = {
          'url': _blog_url,
          'latest_checked': update_time
        }

        updated_members_urls.append(blog_info)
        previous[i]['update_time'] = _update_time

        break
    else:
      blog_info = {
        'url': _blog_url,
        'latest_checked': open_date
      }

      updated_members_urls.append(blog_info)
      new_members.append(_member)

  previous += new_members

  try:
    with open("previous.json", "w") as f:
      json.dump(previous, f, indent=2, ensure_ascii=False)
  except Exception as e:
    raise ValueError(f'failed to write json: {e}')

  return updated_members_urls

def get_new_blog_urls(url, checked):
  latest_checked = datetime.datetime.strptime(checked, "%Y.%m.%d %H:%M")
  urls = []

  for num in count():
    page_url = url + f"&page={num}"
    try:
      soup = fetch_page(page_url)
    except Exception as e:
      raise ValueError(f'failed to fetch_page: {e}')

    for a_tag in soup.find_all('a', attrs={'class', 'hv--thumb'}):
      created = datetime.datetime.strptime(a_tag.find('p', class_='bl--card__date').text, "%Y.%m.%d %H:%M")

      if created <= latest_checked:
        return urls

      blog_url =  base + a_tag.get('href')
      urls.append(blog_url)

    pages = soup.find_all('li', class_='coun')
    if len(pages) == 0:
      break

    last_page = int(pages[-1].text)
    if num+1 >= last_page:
      break

  return urls

def get_blog_info(url):
  try:
    soup = fetch_page(url)
  except Exception as e:
    raise ValueError(f'failed to fetch_page: {e}')

  header = soup.find('header', attrs={'class', 'bd--hd'})
  title = header.find('h1').text
  created = header.find('p').text
  name = soup.find('p', class_='bd--prof__name f--head').text
  content = soup.find('div', attrs={'class', 'bd--edit'}).text

  return {
    'title': title,
    'member': name,
    'created': created,
    'content': content,
    'url': url,
  }

def create_client():
  host = os.environ.get("OPEN_SEARCH_URL")
  user_name = os.environ.get("USER_NAME")
  password = os.environ.get("PASSWORD")

  client =  OpenSearch(
    hosts = [host],
    http_compress = True, # enables gzip compression for request bodies
    http_auth = (user_name, password),
    use_ssl = True,
    verify_certs = False,
    ssl_assert_hostname = False,
    ssl_show_warn = False,
  )

  return client

def index_document(client, data):
  client.index(
    index = 'blogs',
    body = data,
  )

async def handle(request):
  logger.info(msg='request log')
  blog_list_url = "https://www.nogizaka46.com/s/n46/diary/MEMBER"

  try:
    updated_member_urls = get_updated_member_urls(blog_list_url)
  except Exception as e:
    raise ValueError(f'failed to get_updated_member_urls: {e}')

  if len(updated_member_urls) == 0:
    logger.info(msg='no updated!')
    return web.Response(text='no updated!')
  
  try:
    client = create_client()
  except Exception as e:
    raise ValueError(f'failed to create_client: {e}')

  for update_member_url in updated_member_urls:
    try:
      new_blog_urls = get_new_blog_urls(update_member_url['url'], update_member_url['latest_checked'])
    except Exception as e:
      raise ValueError(f'failed to get_new_blog_urls: {e}')

    for new_blog_url in new_blog_urls:
      try:
        data = get_blog_info(new_blog_url)
      except Exception as e:
        raise ValueError(f'failed to get_blog_info: {e}')

      try:
        index_document(client, data)
        logger.info(msg=f'new! {data["title"]}')
      except Exception as e:
        raise ValueError(f'failed to index_document: {e}')

    return web.json_response(new_blog_urls)

app = web.Application()
app.add_routes([web.get('/index', handle)])

if __name__ == "__main__":
  try:
    web.run_app(app, port=14646)
  except Exception as e:
    logger.error(msg=f'Error: {e}')