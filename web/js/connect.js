let username = "";
let password = "";

let refresh_token = "";
let connected = false;

const headers = new Headers();
headers.append("Content-Type", "application/x-www-form-urlencoded");

const getToken = async () => {
	console.log("get token called");

	const params = new URLSearchParams();
	params.append("client_id", "frontend_client");

	const requestOptions = {
		method: "POST",
		headers: headers,
		body: params,
	};

	if (refresh_token == "") {
		params.append("username", username);
		params.append("password", password);
		params.append("grant_type", "password");
	} else {
		params.append("refresh_token", refresh_token);
		params.append("grant_type", "refresh_token");
	}
	console.log(requestOptions.body)
	return fetch(
		"http://127.0.0.1:8090/realms/sirius/protocol/openid-connect/token",
		requestOptions
	)
		.then((response) => response.json())
		.then((result) => {
			console.log("token accepted");
			refresh_token = result.refresh_token;

			return result.access_token;
		})
		.catch((error) => console.error(error));
};

const client = new Centrifuge(
	"ws://127.0.0.1:8080/centrifugo/connection/websocket",
	{
		getToken: getToken,
		debug: true,
	}
);

client.on("connected", () => {
	password = null;

	document.getElementById("password").value = "";

	connected = true;
});

client.on("disconnected", () => {
	refresh_token = "";
	client.setToken("");

	connected = false;
});

document.getElementById("btnLogin").addEventListener("click", () => {
	if (connected) {
		return;
	}

	username = document.getElementById("username").value;
	password = document.getElementById("password").value;

	if (username != "" && password != "") {
		client.connect();
	} else {
		console.error("empty username or password");
	}
});

document.getElementById("btnLogout").addEventListener("click", () => {
	if (!connected) {
		return;
	}

	console.log("disconnecting");

	const params = new URLSearchParams();
	params.append("client_id", "frontend_client");

	const requestOptions = {
		method: "POST",
		headers: headers,
		body: params,
	};

	if (refresh_token != "") {
		params.append("refresh_token", refresh_token);
		params.append("grant_type", "refresh_token");

		fetch(
			"http://127.0.0.1:8090/realms/sirius/protocol/openid-connect/logout",
			requestOptions
		)
			.then((response) => {
				console.log("session closed");
			})
			.catch((error) => console.error(error))
			.finally(() => {
				client.disconnect();
			});
	} else {
		console.error("no refresh token");
	}
});

document.getElementById("btnSubscribe").addEventListener("click", () => {
    const channelName = document.getElementById("channel").value;
    
    if (!channelName) {
        console.error("Channel name is required");
        return;
    }

    const sub = client.subscribe(channelName, {
        getToken: getToken
    });

    sub.on("publication", function(ctx) {
        const output = document.getElementById("output");
        output.innerHTML += `<p>New message: ${JSON.stringify(ctx.data)}</p>`;
    });

    sub.on("subscribing", function(ctx) {
        console.log(`Subscribing to ${channelName}`);
    });

    sub.on("subscribed", function(ctx) {
        console.log(`Subscribed to ${channelName}`);
    });

    sub.on("error", function(ctx) {
        console.error(`Subscription error: ${ctx.error}`);
    });
});

document.getElementById("btnSend").addEventListener("click", () => {
    const channelName = document.getElementById("channel").value;
    const messageText = document.getElementById("message").value;
    
    if (!channelName || !messageText) {
        console.error("Channel and message are required");
        return;
    }

    client.publish(channelName, {
        text: messageText
    }).then(function(res) {
        console.log("Message published", res);
    }, function(err) {
        console.error("Publish error", err);
    });
});