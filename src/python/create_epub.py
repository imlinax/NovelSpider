#!/usr/bin/python2
# -*- coding: utf-8 -*-
'''
create novel epub
'''

from ebooklib import epub
import os

def txt2html(path):
    html = ""
    f = open(path)
    for line in f.readlines():
        html += '<p>' + line + '</p>'

    f.close()
    return html

book = epub.EpubBook()
book.add_author("chen dong")



chapters = []
rootdir = u'/Users/wangli/code/novel_spider/golang/圣墟'
list = os.listdir(rootdir) #列出文件夹下所有的目录与文件
list.sort()
for i in range(0,len(list)):
    path = os.path.join(rootdir,list[i])
    if os.path.isfile(path):
        c1 = epub.EpubHtml(title=os.path.basename(path),file_name="chap%d.xhtml" % i,lang='zh')
        c1.content = txt2html(path)
        chapters.append(c1)
        book.add_item(c1)

book.set_identifier("123456")
book.set_title(u"圣墟")
book.set_language('en')

c1 = epub.EpubHtml(title='Introduction', file_name='intro.xhtml', lang='en')
    c1.content=u'<html><head></head><body><h1>圣墟</h1><p>Introduction paragraph where i explain what is happening.</p></body></html>'

book.toc = (epub.Link('intro.xhtml', 'Introduction', 'intro'),
             (epub.Section('Simple book'),
             tuple(chapters))
            )

# add default NCX and Nav file
book.add_item(epub.EpubNcx())
book.add_item(epub.EpubNav())

# define CSS style
style = 'BODY {color: white;}'
nav_css = epub.EpubItem(uid="style_nav", file_name="style/nav.css", media_type="text/css", content=style)

# add CSS file
book.add_item(nav_css)

# basic spine
lst = [['nav'], chapters]
book.spine = sum(lst, [])


epub.write_epub("sx.epub", book, {})