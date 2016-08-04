/*

Package htmltester includes the tester.html and makes it available under the /tester.html path.

Include this packages if you want the executable binary to self-contain the tester page.

*/
package htmltester

import (
	"io"
	"net/http"
)

func init() {
	http.HandleFunc("/tester.html", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, tester_html)
	})
}

var tester_html = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>productws demo tester</title>
<style>
	html {display: table; margin: auto;}
	body {display: table-cell; font-family: Arial; color: #666;}
	textarea {height: 60px; color: #222;}
	code {color: green; background: #eee; padding: 3px;}
	h1 {text-align:center; color:#222;}
	h2 {font-size:1.3em; margin:10px 0px 10px 0px; color:black; border-bottom: 1px solid #aa8; background-color: #fff8d0; padding: 3px;}
	button {width: 120px; margin:3px 0px 3px 0px; display: inline-block; color: white; background: #007aec; border-radius: 3px;}
	button:hover {background: #005aca; text-decoration: none;}
	.req, .resp {padding-left: 5px; width: 750px;}
	.req {border-left: 10px solid #9e9;}
	.resp {border-left: 10px solid #faa; background-color: #f6f6f6;}
	.footer {background: #fff; padding: 6px 10px 6px 10px; border-top: 1px solid #aaa; font-style: italic; font-size: 0.8em; margin-top:3px; text-align:center;}
</style>
</head>

<body>
	<h1>productws demo tester</h1>
	<div id="content">
		Loading tester application...<br> If this text does not
		disappear shortly, the React app has problems startup up.
	</div>

	<script src="https://npmcdn.com/react@15.3.0/dist/react.js"></script>
	<script src="https://npmcdn.com/react-dom@15.3.0/dist/react-dom.js"></script>
	<script src="https://npmcdn.com/babel-core@5.8.38/browser.min.js"></script>
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
	
	<script type="text/babel">

// Enable cross domain request (IE, Edge needs it!)
$.support.cors = true;

function errorObj(xhr, status, err) {
	return {response: "ERROR: " + xhr.status + " " + err.toString() + "\n" + xhr.responseText
		+(xhr.readyState == 0 ? "\nService unavailable! Wrong Service URL?" : "")};
}

function serviceURL() {
	return document.getElementById("serviceURLId").value;
}

var ServiceURL = React.createClass({
	render: function() {
		return <div className="serviceURL">
			Service URL:<br/>
			<input className="req" id="serviceURLId" type="text" defaultValue="http://localhost:8081/" />
		</div>;
	}
});

var CreateCall = React.createClass({
	getInitialState: function () {
		return { request: '{"Name":"JSCO Mouse","Desc":"Computer Optical Noiseless Mouse","Prices":{"USD":{"Value":2782,"Multiplier":100}}}' };
	},
	render: function() {
		return <div>
			Product JSON:<br/>
			<textarea className="req" value={this.state.request} onChange={this.onChange}></textarea><br/>
			<button onClick={this.handleClick}>CREATE</button> <code>"POST /create"</code><br/>
			Response:<br/>
			<textarea className="resp" readOnly="true" value={this.state.response}></textarea>
		</div>;
	},
	onChange: function(event) {
		this.setState({request: event.target.value});
	},
	handleClick: function() {
		var t = this;
		t.setState({response : "Calling web service..."});
		$.ajax({
			url: serviceURL() + "create",
			type: 'POST',
			data: t.state.request,
			dataType: 'json',
			success: function(data) { t.setState({response : JSON.stringify(data)}); },
			error: function(xhr, status, err) {	t.setState(errorObj(xhr, status, err)); }
		});
	}
});

var ListCall = React.createClass({
	getInitialState: function () {
		return {};
	},
	render: function() {
		return <div>
			<button onClick={this.handleClick}>LIST</button> <code>"GET /list"</code><br/>
			Response:<br/>
			<textarea className="resp" readOnly="true" value={this.state.response}></textarea>
		</div>;
	},
	handleClick: function() {
		var t = this;
		t.setState({response : "Calling web service..."});
		$.ajax({
			url: serviceURL() + "list",
			success: function(data) { t.setState({response : JSON.stringify(data)}); },
			error: function(xhr, status, err) {	t.setState(errorObj(xhr, status, err)); }
		});
	}
});

var DetailsCall = React.createClass({
	getInitialState: function () {
		return {id: "3"};
	},
	render: function() {
		return <div>
			Product ID:<br/>
			<input className="req" type="text" value={this.state.id} onChange={this.onChange}></input><br/>
			<button onClick={this.handleClick}>DETAILS</button> <code>"GET /details/{this.state.id}"</code><br/>
			Response:<br/>
			<textarea className="resp" readOnly="true" value={this.state.response}></textarea>
		</div>;
	},
	onChange: function(event) {
		this.setState({id: event.target.value});
	},
	handleClick: function() {
		var t = this;
		t.setState({response : "Calling web service..."});
		$.ajax({
			url: serviceURL() + "details/" + t.state.id,
			success: function(data) { t.setState({response : JSON.stringify(data)}); },
			error: function(xhr, status, err) {	t.setState(errorObj(xhr, status, err)); }
		});
	}
});

var UpdateCall = React.createClass({
	getInitialState: function () {
		return { request: '{"ID":3,"Name":"JSCO Mouse","Desc":"Computer Optical Noiseless Mouse","Tags":["Computer","Mouse"],"Prices":{"USD":{"Value":2782,"Multiplier":100},"GBP":{"Value":2093,"Multiplier":100}}}' };
	},
	render: function() {
		return <div>
			Product JSON:<br/>
			<textarea className="req" value={this.state.request} onChange={this.onChange}></textarea><br/>
			<button onClick={this.handleClick}>UPDATE</button> <code>"PUT /update"</code><br/>
			Response:<br/>
			<textarea className="resp" readOnly="true" value={this.state.response}></textarea>
		</div>;
	},
	onChange: function(event) {
		this.setState({request: event.target.value});
	},
	handleClick: function() {
		var t = this;
		t.setState({response : "Calling web service..."});
		$.ajax({
			url: serviceURL() + "update",
			type: 'PUT',
			data: t.state.request,
			dataType: 'json',
			success: function(data) { t.setState({response : JSON.stringify(data)}); },
			error: function(xhr, status, err) {	t.setState(errorObj(xhr, status, err)); }
		});
	}
});

var SetpricesCall = React.createClass({
	getInitialState: function () {
		return { request: '{"ID":3,"Prices":{"GBP":{"Value":1999,"Multiplier":100},"HUF":{"Value":7717,"Multiplier":1}}}' };
	},
	render: function() {
		return <div>
			Product JSON:<br/>
			<textarea className="req" value={this.state.request} onChange={this.onChange}></textarea><br/>
			<button onClick={this.handleClick}>SET PRICES</button> <code>"PUT /setprices"</code><br/>
			Response:<br/>
			<textarea className="resp" readOnly="true" value={this.state.response}></textarea>
		</div>;
	},
	onChange: function(event) {
		this.setState({request: event.target.value});
	},
	handleClick: function() {
		var t = this;
		t.setState({response : "Calling web service..."});
		$.ajax({
			url: serviceURL() + "setprices",
			type: 'PUT',
			data: t.state.request,
			dataType: 'json',
			success: function(data) { t.setState({response : JSON.stringify(data)}); },
			error: function(xhr, status, err) {	t.setState(errorObj(xhr, status, err)); }
		});
	}
});

var TesterApp = React.createClass({
	render: function() {
		return <div>
			<ServiceURL/>
			<h2>Create a new product</h2> <CreateCall/>
			<h2>List all product IDs</h2> <ListCall/>
			<h2>Get details of a product</h2> <DetailsCall/>
			<h2>Update a product</h2> <UpdateCall/>
			<h2>Set price points</h2> <SetpricesCall/>
			<div className="footer">Visit <a href="https://github.com/icza/productws" target="_blank">github.com/icza/productws</a></div>
		</div>;
	}
});

ReactDOM.render(
	<TesterApp />,
	document.getElementById('content')
);

</script>

</body>
</html>`
