import requests
from bs4 import BeautifulSoup
import json
import re
import datetime
from itertools import count
import os
from opensearchpy import OpenSearch

base = 'https://www.nogizaka46.com'
open_date = '2011.11.11 00:00'

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

def get_updated_member_urls(url):
  soup = BeautifulSoup(requests.get(url).text, "html.parser")
  if soup == "":
    return

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

  with open("previous.json", "r") as f:
    previous = json.loads(f.read())

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

  with open("previous.json", "w") as f:
    json.dump(previous, f, indent=2, ensure_ascii=False)

  return updated_members_urls

def get_new_blog_urls(url, checked):
  latest_checked = datetime.datetime.strptime(checked, "%Y.%m.%d %H:%M")
  urls = []

  for num in count():
    page_url = url + f"&page={num}"
    soup = BeautifulSoup(requests.get(page_url).text, "html.parser")
    if soup == "":
      return

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
  soup = BeautifulSoup(requests.get(url).text, "html.parser")
  if soup == "":
    return

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

  return OpenSearch(
    hosts = [host],
    http_compress = True, # enables gzip compression for request bodies
    http_auth = (user_name, password),
    use_ssl = True,
    verify_certs = False,
    ssl_assert_hostname = False,
    ssl_show_warn = False,
)

def index_document(client, data):
  client.index(
    index = 'blogs',
    body = data,
  )

if __name__ == "__main__":
  blog_list_url = "https://www.nogizaka46.com/s/n46/diary/MEMBER"
  updated_member_urls = get_updated_member_urls(blog_list_url)

  if len(updated_member_urls) != 0:
    client = create_client()

    for update_member_url in updated_member_urls:
      new_blog_urls = get_new_blog_urls(update_member_url['url'], update_member_url['latest_checked'])
      for new_blog_url in new_blog_urls:
        data = get_blog_info(new_blog_url)
        index_document(client, data)