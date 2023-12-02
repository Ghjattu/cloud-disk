const baseURL = "http://127.0.0.1:8082";

const CheckFileExistence = async (fileHash) => {
	const url = `${baseURL}/file/exist/${fileHash}`;
	const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE1ODI1ODUsImlhdCI6MTcwMTQ5NjE4NSwibmFtZSI6InRlc3QiLCJ1c2VyX2lkIjoxfQ.g9A9DuvWEoVgFlxV-jk4jOn_R7MFzBKghKgtMHsEr5Y';
	const headers = {
		'Content-Type': 'application/json',
		'Accept': 'application/json',
		'Authorization': `Bearer ${token}`,
	};

	const response = await fetch(url, {
		method: 'GET',
		headers: headers,
	});
	return await response.json();
};

const UploadFileAPI = {
	CheckFileExistence
};

export default UploadFileAPI;