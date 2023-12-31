import React, { useEffect } from 'react';
import './Dashboard.css';
import PropTypes from 'prop-types';
import { useState } from 'react';
import { useImmer } from 'use-immer';
import CalculateMD5 from '../utils/calculateMD5.js';
import uploadFileAPI from '../api/uploadFileAPI.js';
import GetFileAPI from '../api/getFilesAPI.js';
import FileList from '../FileList/FileList.js';
import LinearDeterminate from '../LinearDeterminate/LinearDeterminate.js';

const Dashboard = ({ token }) => {
	const [selectedFile, setSelectedFile] = useState(null);
	const [isUploading, setIsUploading] = useState(false);
	const [fileList, setFileList] = useImmer([]);
	const [progress, setProgress] = useState(0);
	const controller = new AbortController();

	const addProgress = (value) => {
		setProgress(oldProgress => {
			return Math.min(100, oldProgress + value);
		});
	};

	useEffect(() => {
		const initiateFileList = async () => {
			const resp = await GetFileAPI.GetFileList(token);
			setFileList(resp.data);
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

		setIsUploading(true);

		// calculate MD5 hash of file
		let fileHash = null;
		let chunksHash = [];
		try {
			const resp = await CalculateMD5(selectedFile);
			fileHash = resp.fileHash;
			chunksHash = resp.chunksHash;
		} catch (error) {
			console.log('Error calculating MD5: ', error);
			setIsUploading(false);
			return;
		}

		// check if file already exists
		const resp = await uploadFileAPI.CheckFileExistence(fileHash, token);
		if (resp.data.exist) {
			alert('File uploaded successfully');
			setIsUploading(false);
			return;
		}
		const uploadedChunksHash = resp.data.chunks_hash;

		// upload file in chunks
		try {
			await uploadFileAPI.
				UploadChunks(selectedFile, fileHash, chunksHash, uploadedChunksHash, token, addProgress, controller);
		} catch (error) {
			alert(error);
			setIsUploading(false);
			return;
		}

		// merge chunks
		try {
			const resp = await uploadFileAPI.MergeChunks(fileHash, selectedFile.name, selectedFile.size, token);
			setFileList(draft => {
				draft.push({
					'file_id': resp.file_id,
					'file_name': selectedFile.name,
					'file_size': selectedFile.size,
					'file_url': resp.file_url,
					'upload_time': resp.upload_time,
				});
			});
			alert('File uploaded successfully');
		} catch (error) {
			alert(error);
		}

		setIsUploading(false);
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
						<button type='submit' disabled={isUploading}>Submit</button>
					</div>
				</form>
			</div>
			<div className='upload-progress'>
				<LinearDeterminate progress={progress}></LinearDeterminate>
			</div>
			<FileList fileList={fileList}></FileList>
		</div>
	);
};

Dashboard.propTypes = {
	token: PropTypes.string.isRequired,
};

export default Dashboard;