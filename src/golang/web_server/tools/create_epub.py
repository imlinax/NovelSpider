#!/usr/bin/python
# -*- coding: utf-8 -*-
'''
create novel epub
'''

from ebooklib import epub
from optparse import OptionParser
import os

def txt2html(path):
    html = ""
    f = open(path)
    for line in f.readlines():
        html += '<p>' + line + '</p>'

    f.close()
    return html

def create_epub(input_dir):
    book = epub.EpubBook()
    book.add_author("unkown")
    
    chapters = []
    rootdir = input_dir
    title=os.path.basename(rootdir)
    list = os.listdir(rootdir) #列出文件夹下所有的目录与文件
    list.sort()
    for i in range(0,len(list)):
        path = os.path.join(rootdir,list[i])
        if os.path.isfile(path):
            c1 = epub.EpubHtml(title=unicode(os.path.basename(path),"utf-8"),file_name="chap%d.xhtml" % i,lang='zh')
            c1.content = txt2html(path)
            chapters.append(c1)
            book.add_item(c1)

    book.set_identifier("123456")
    book.set_title(unicode(title,"utf-8"))
    book.set_language('zh_CN')

    c1 = epub.EpubHtml(title='Introduction', file_name='intro.xhtml', lang='zh_CN')
    c1.content=u'<html><head></head><body><h1>%s</h1><p>Introduction paragraph where i explain what is happening.</p></body></html>' % unicode(title,"utf-8")

    #book.toc = (epub.Link('intro.xhtml', 'Introduction', 'intro'),
    #             (epub.Section('Simple book'),
    #             tuple(chapters))
    #            )

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


    epub.write_epub("%s.epub" % title, book, {})

def get_user_paras():
    try:
        opt = OptionParser()
        opt.add_option('--dir',
        dest='input_dir',
        type=str,
        help='where the novel store in (input)')

        (options, args) = opt.parse_args()
        return options

    except Exception as ex:
        print('exception :{0}'.format(str(ex)))
        return None

def main():
    user_paras = get_user_paras()
    create_epub(user_paras.input_dir)
if __name__ == "__main__":
    main()