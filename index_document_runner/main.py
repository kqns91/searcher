import requests
from bs4 import BeautifulSoup
from dotenv import load_dotenv
import json
import re
import datetime
from itertools import count

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

  for _member in members:
    _name = _member['name']
    _update_time = _member['update_time']
    _blog_url = _member['blog_url']

    for member in previous:
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
        break
    else:
      blog_info = {
        'url': _blog_url,
        'latest_checked': open_date
      }

      updated_members_urls.append(blog_info)

  with open("previous.json", "w") as f:
    json.dump(members, f, indent=2, ensure_ascii=False)

  return updated_members_urls

if __name__ == "__main__":
  blog_list_url = "https://www.nogizaka46.com/s/n46/diary/MEMBER"
  updated_member_urls = get_updated_member_urls(blog_list_url)