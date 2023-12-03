import React from 'react';
import { useState } from 'react';
import './App.css';
import CalculateMD5 from './utils/calculateMD5.js';
import uploadFileAPI from './api/uploadFileAPI.js';

const App = () => {
	const [selectedFile, setSelectedFile] = useState(null);

	const handleFileChange = (event) => {
		setSelectedFile(event.target.files[0]);
	};

	const handleSubmit = async (event) => {
		event.preventDefault();

		if (!selectedFile) {
			alert('No file selected');
			return;
		}

		// calculate MD5 hash of file
		let fileHash = null;
		let chunksHash = [];
		try {
			const resp = await CalculateMD5(selectedFile);
			fileHash = resp.fileHash;
			chunksHash = resp.chunksHash;
		} catch (error) {
			console.log('Error calculating MD5: ', error);
			return;
		}

		// check if file already exists
		const resp = await uploadFileAPI.CheckFileExistence(fileHash);
		if (resp.exist) {
			alert('File uploaded successfully');
			return;
		}
		const uploadedChunksHash = resp.chunks_hash;

		// upload file in chunks
		try {
			const resp = await uploadFileAPI.
				UploadFileInChunks(selectedFile, fileHash, chunksHash, uploadedChunksHash);
			console.log(resp);
		} catch (error) {
			console.log('Error uploading file: ', error);
			return;
		}
	};

	return (
		<div className='form-container'>
			<form onSubmit={handleSubmit}>
				<div>
					<label htmlFor='file'>Select a file</label>
				</div>
				<div>
					<input type='file' id='file' name='file' onChange={handleFileChange} />
				</div>
				<div>
					<button type='submit'>Submit</button>
				</div>
			</form>
		</div>
	);
};

export default App;