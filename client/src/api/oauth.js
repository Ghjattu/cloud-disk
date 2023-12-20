const GetGithubLoginURL = async () => {
	const url = '/oauth/github/login_url';

	const response = await fetch(url, {
		method: 'GET',
		headers: {
			'Accept': 'application/json',
		},
	});

	return response.json();
};

const LoginWithGithub = async (token) => {
	const url = `/oauth/github/callback/${token}`;

	const response = await fetch(url, {
		method: 'POST',
		headers: {
			'Accept': 'application/json',
		},
	});

	return response.json();
};

const OAuthAPI = {
	GetGithubLoginURL,
	LoginWithGithub,
};

export default OAuthAPI;