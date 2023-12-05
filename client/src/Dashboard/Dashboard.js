import React, { useEffect } from 'react';
import './Dashboard.css';
import PropTypes from 'prop-types';
import { useState } from 'react';
import { useImmer } from 'use-immer';
import CalculateMD5 from '../utils/calculateMD5.js';
import uploadFileAPI from '../api/uploadFileAPI.js';
import GetFileAPI from '../api/getFilesAPI.js';
import FormatBytes from '../utils/formatBytes.js';

const Dashboard = ({ token }) => {
	const [selectedFile, setSelectedFile] = useState(null);
	const [fileList, setFileList] = useImmer([]);

	useEffect(() => {
		const initiateFileList = async () => {
			const resp = await GetFileAPI.GetFileList(token);
			if (resp != null) {
				setFileList(resp);
			}
		};
		initiateFileList();
	}, []);

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
		const resp = await uploadFileAPI.CheckFileExistence(fileHash, token);
		if (resp.exist) {
			alert('File uploaded successfully');
			return;
		}
		const uploadedChunksHash = resp.chunks_hash;

		// upload file in chunks
		try {
			const resp = await uploadFileAPI.
				UploadFileInChunks(selectedFile, fileHash, chunksHash, uploadedChunksHash, token);
			if (resp.file_success) {
				setFileList((draft) => {
					draft.push({
						file_id: resp.file_id,
						file_name: selectedFile.name,
						file_size: selectedFile.size,
						file_url: resp.file_url,
					});
				});
				setSelectedFile(null);
				alert('File uploaded successfully');
			} else {
				alert('Error uploading file');
			}
		} catch (error) {
			alert('Error uploading file: ', error);
			return;
		}
	};

	return (
		<div className='dashboard-container'>
			<div className='form-container'>
				<form onSubmit={handleSubmit}>
					<div>
						<label htmlFor='file'>Upload a file</label>
					</div>
					<div>
						<input type='file' id='file' name='file' onChange={handleFileChange} />
					</div>
					<div>
						<button type='submit'>Submit</button>
					</div>
				</form>
			</div>
			<div className='file-list-container'>
				<h3 className='file-list-title'>File List</h3>
				<ul className='file-list'>
					{fileList.map((file) => (
						<li key={file.file_id} className='file-list-item'>
							<span className='file-name'>{file.file_name}</span>
							<span className='file-size'>{ FormatBytes(file.file_size) }</span>
							<span className='file-url'><a href={file.file_url}>Download</a></span>
						</li>
					))}
				</ul>
			</div>
		</div>
	);
};

Dashboard.propTypes = {
	token: PropTypes.string.isRequired,
};

export default Dashboard;