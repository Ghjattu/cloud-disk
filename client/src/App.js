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
		try {
			fileHash = await CalculateMD5(selectedFile);
			console.log(fileHash);
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

		// TODO: upload file in chunks

		// const formData = new FormData();
		// formData.append('file', selectedFile);
		// fetch('http://localhost:5000/upload', {
		// 	method: 'POST',
		// 	body: formData,
		// })
		// 	.then((response) => response.json())
		// 	.then((result) => {
		// 		console.log('Success:', result);
		// 	})
		// 	.catch((error) => {
		// 		console.error('Error:', error);
		// 	});
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