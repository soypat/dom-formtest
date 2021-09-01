package main

// UI events.
const (
	// The load event fires when the webpage finishes loading.
	// It can also fire on nodes of elements like images, scripts, or objects.
	EventLoad = "load"
	// Fires when a user leaves the page. Usually because they click on a link.
	EventUnload = "unload"
	// This event fires when the browser encounters a JavaScript error or an asset that doesn’t exist.
	EventError = "error"
	// It fires when we resize the browser window. But browsers repeatedly fire this event,
	// so avoid using this event to trigger complicated code; it might make the page less responsive.
	EventResize = "resize"
	// This event fires when the user scrolls up/down on the browser window.
	// It can relate to the entire page or a specific element on the page.
	EventScroll = "scroll"
)

// Focus and blur events.
const (
	// This event fires, for a specific DOM node, when an element gains focus.
	EventFocus = "focus"
	// This fires, for a specific DOM node, when an element loses focus.
	EventBlur = "blur"
)

// Human interface device input events.
const (
	// Mouse events or touchpad
	// This event fires when the user clicks on the primary mouse button (usually the left button).
	// This event also fires if the user presses the Enter key on the keyboard when an element has focus.
	// Touch-screen: A tap on the screen acts like a single primary mouse button click.
	EventClick = "click"
	// It fires when the user clicks down on any mouse button.
	EventMouseDown, EventTouchstart = "mousedown", "touchstart"

	// We have separate mousedown and mouseup events to add drag-and-drop functionality
	// or controls in game development. Don’t forget a click event is the combination
	// of mousedown and mouseup events.

	// It fires when the user releases a mouse button.
	EventMouseUp = "mouseup"
	// It fires when the user moves the cursor, which was inside an element before, outside the element.
	// We can say that it fires when the cursor moves off the element.
	EventMouseOut = "mouseout"
	// It fires when the user moves the cursor, which was outside an element before, inside the element.
	// We can say that it fires when we move the cursor over the element.
	EventMouseOver = "mouseover"
	// It fires when the user moves the cursor around the element. This event is frequently triggered.
	EventMouseMove = "mousemove"

	// Keyboard events

	// The keydown and keypress events fire before a character appears on the screen,
	// the keyup fires after it shows. To know the key pressed when you use the
	// keydown and keypress events, the event object has a keyCode property.
	// This property, instead of returning the letter for that key, returns
	// the ASCII code of the lowercase for that key.

	// The keyup event fires when the user releases a key on the keyboard.
	EventKeyUp = "keyup"
	// It fires when the user presses any key in the keyboard.
	// If the user holds down the key, this event fires repeatedly.
	EventKeyDown = "keydown"
	// It fires when the user presses a key that results in printing a character on the screen.
	// This event fires repeatedly if the user holds down the key.
	// This event will not fire for the enter, tab, or arrow keys; the keydown event would.
	EventKeyPress = "keypress"

	// This event fires when the user clicks the primary mouse button, in quick succession, twice.
	EventDoubleClick = "dblclick"
)

// Form events
const (
	// This event fires on the node representing the <form> element when a user submits a form.
	EventSubmit = "submit"
	// It fires when the status of various form elements change.
	// This is a better option than using the click event because clicking is not the only way users interact with the form.
	EventChange = "change"
	// This event fires when the value of an <input> or a <textarea> changes
	// (doesn’t fire for deleting in IE9). You can use keydown as a fallback in older browsers.
	EventInput = "input"
)

// HTML5 events.
const (
	// This event triggers when the DOM tree forms i.e. the script is loading.
	// Scripts start to run before all the resources like images, CSS, and JavaScript loads.
	// You can attach this event either to the window or the document objects.
	EventDOMContentLoaded = "DOMContentLoaded"
	// It fires when the URL hash changes without refreshing the entire window.
	// Hashes (#) link specific parts (known as anchors) within a page.
	//  It works on the window object; the event object contains both the oldURL
	// and the newURL properties holding the URLs before and after the hashchange.
	EventHashChange = "hashchange"
	// This event fires on the window object just before the page unloads.
	// This event should only be helpful for the user,
	// not encouraging them to stay on the page.
	// You can add a dialog box to your event, showing a message alerting the users like their changes are not saved.
	EventBeforeUnload = "beforeunload"
)
