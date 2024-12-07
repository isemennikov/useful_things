#!/usr/bin/env python
from typing import Union, Optional, Awaitable

import tornado.web
import tornado.websocket
import tornado.ioloop

class WebSocketHandler(tornado.websocket.WebSocketHandler):
    def open(self):
        print("WebSocket opened")

    def on_message(self, message):
        self.write(f"You said: {message}")

    def on_close(self):
        print("WebSocket closed")

app = tornado.web.Application([(r'/ws', WebSocketHandler)])

if __name__=="__main__":
    app.listen(8765)
    tornado.ioloop.IOLoop.current().start()
