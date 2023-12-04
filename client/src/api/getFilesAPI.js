const baseURL = 'http://127.0.0.1:8082';

const GetFileList = async (token) => {
	const url = `${baseURL}/file/list`;
	const headers = {
		'Accept': 'application/json',
		'Authorization': `Bearer ${token}`,
	};

	const response = await fetch(url, {
		method: 'GET',
		headers: headers,
	});

	return await response.json();
};

const GetFileAPI = {
	GetFileList,
};

export default GetFileAPI;