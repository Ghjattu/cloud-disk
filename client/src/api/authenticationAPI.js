const baseURL = 'http://127.0.0.1:8081';

const Register = async (username, password) => {
	const url = `${baseURL}/user/register`;
	const headers = {
		'Content-Type': 'application/json',
		'Accept': 'application/json',
	};

	const response = await fetch(url, {
		method: 'POST',
		headers: headers,
		body: JSON.stringify({
			'name': username,
			'password': password,
		}),
	});

	return await response.json();
};

const Login = async (username, password) => {
	const url = `${baseURL}/user/login`;
	const headers = {
		'Content-Type': 'application/json',
		'Accept': 'application/json',
	};

	const response = await fetch(url, {
		method: 'POST',
		headers: headers,
		body: JSON.stringify({
			'name': username,
			'password': password,
		}),
	});

	return await response.json();
};

const AuthenticationAPI = {
	Register,
	Login,
};

export default AuthenticationAPI;