Выкачать файлик с веба по адресу.

На вход в stdin json-объект с полями:
url, follow-redirects, timeout, ignore-ssl-certificate, output

На выходе в stdout должен выдать json с:
success, http-code, content-length, content-type, redirects[]

И еще на указанный в output путь по tcp выстримить скачанное

