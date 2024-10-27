# This is a simple example to using htmx and ws.

The actual output is trivial, it pushes a time every three seconds to the subscribers.

However it is an evolution of the standard boiler plate that I've been using as it includes a number of newer features:

	middleware sessionId capture, no need to call get cookies at the start of each method
		-- doesn't yet return a page from the page array. I can't decide if that would be good
		or not

	getHtmlFromHTML method included with related structs

	websocket boiler included.

	