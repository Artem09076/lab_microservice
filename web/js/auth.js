const kc = new Keycloak({
	url: "http://127.0.0.1:8090",
	realm: "sirius",
	clientId: "frontend_client",
});

const initAuth = async () => {
	kc.init({
		onLoad: "login-required",
	});

	kc.onReady = (auth) => {
		console.log("KC: Ready:", auth);
		document.body.classList.remove("hidden");

		document.getElementById("btnLogout").addEventListener("click", () => {
			kc.logout();
		});

		console.log("token", kc.token);
	};

	kc.onAuthSuccess = () => console.log("KC: AuthSuccess:");

	kc.onAuthError = () => console.log("KC: AuthError:");

	kc.onAuthRefreshSuccess = () => console.log("KC: AuthRefreshSuccess:");

	kc.onAuthRefreshError = () => console.log("KC: AuthRefreshError:");

	kc.onAuthLogout = () => console.log("KC: AuthLogout:");

	kc.onTokenExpired = () => {
		console.log("KC: TokenExpired:");
		kc.updateToken();
	};
};

initAuth();
