import GetChunkSize from '../utils/getChunkSize.js';

const baseURL = 'http://127.0.0.1:8082';
const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE2Njk0OTUsImlhdCI6MTcwMTU4MzA5NSwibmFtZSI6InRlc3QiLCJ1c2VyX2lkIjoxfQ.t1UBhQ_yMPU-u6JwhUlo3UVLAvbH_DRjWx9iHHI4ptY';

const CheckFileExistence = async (fileHash) => {
	const url = `${baseURL}/file/exist/${fileHash}`;
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

const UploadFileInChunks = async (file, fileHash, chunksHash, uploadedChunksHash) => {
	console.log('start uploading file in chunks');
	const url = `${baseURL}/file/upload_chunk`;
	const headers = {
		// add this line will cause error
		// 'Content-Type': 'multipart/form-data',
		'Authorization': `Bearer ${token}`,
	};

	// TODO: need refactor
	return new Promise((resolve, reject) => {
		const chunkSize = GetChunkSize(); // 100 KB
		const totalChunks = Math.ceil(file.size / chunkSize);

		for (let chunkNum = 0; chunkNum < totalChunks; chunkNum++) {
			if (uploadedChunksHash.includes(chunksHash[chunkNum])) {
				continue;
			}

			const start = chunkNum * chunkSize;
			const end = ((start + chunkSize) >= file.size) ? file.size : start + chunkSize;
			const chunk = file.slice(start, end);

			const formData = new FormData();
			formData.append('chunk', chunk);
			formData.append('chunk_info', JSON.stringify({
				'file_hash': fileHash,
				'total_chunks': totalChunks,
				'chunk_hash': chunksHash[chunkNum],
			}));

			fetch(url, {
				method: 'POST',
				headers: headers,
				body: formData,
			})
				.then((response) => {
					if (!response.ok) {
						reject(new Error('Network response was not ok'));
					}
					return response.json();
				})
				.then((resp) => {
					console.log('Upload chunk %d successfully', chunkNum);
					if (resp.file_success) {
						resolve(resp);
					}
					if (!resp.chunk_success) {
						reject(new Error('Upload chunk failed'));
					}
				})
				.catch((error) => {
					reject(error);
				});
		}

		// TODO: delete this line
		resolve('all fetch done');
	});
};

const UploadFileAPI = {
	CheckFileExistence,
	UploadFileInChunks,
};

export default UploadFileAPI;