const GetFileList = async (token) => {
	const url = 'file/list';
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