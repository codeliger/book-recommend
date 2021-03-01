import requests, json, random

genres = ['fiction', 'nonfiction', 'scifi', 'romance', 'autobiography', 'travel', 'action', 'adventure', 'classic', 'mystery', 'detective', 'fantasy', 'horror', 'thriller', 'cooking', 'modern', 'classic']

books = []

id_count = 1

for genre in genres:
  response = requests.get('https://www.googleapis.com/books/v1/volumes?q=subject:{}'.format(genre))
  if not response.ok:
    raise RuntimeError("failed to fetch books for genre {}: {}".format(genre, response.text()))
  
  rawbooks = response.json()['items']

  if len(rawbooks) == 0:
    print('genre', genre, 'returned no results')
  for rawbook in rawbooks:
    try:
      book = {
        "id": id_count
      }
      book['title'] = rawbook['volumeInfo']['title']
      book['authors'] = rawbook['volumeInfo']['authors']
      book['genres'] = [genre]
      if 'pageCount' in rawbook['volumeInfo']:
        book['pages'] = rawbook['volumeInfo']['pageCount']
      else:
        book['pages'] = random.randint(50,500)
      book['publication_year'] = random.randint(1930,2010)
      book['rating'] = random.randint(1,5)
      books.append(book)
      id_count += 1
    except KeyError as e:
      print("book is missing", e.args[0])
      continue


with open('books.json', 'w') as f:
  json.dump(books, f)
    