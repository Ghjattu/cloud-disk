import React from 'react';
import './FileList.css';
import PropTypes from 'prop-types';
import FormatBytes from '../utils/formatBytes.js';
import FormatTime from '../utils/formatTime.js';

const FileList = ({ fileList }) => {
	const handleDownload = (e, filePath, fileName) => {
		e.preventDefault();
		const url = `/static/videos/${filePath}`;

		fetch(url)
			.then(response => response.blob())
			.then(blob => {
				const a = document.createElement('a');
				const url = URL.createObjectURL(blob);
				a.href = url;
				a.download = fileName;
				document.body.appendChild(a);
				a.click();
				document.body.removeChild(a);
				URL.revokeObjectURL(url);
			})
			.catch(error => {
				console.log(error);
			});
	};

	return (
		<div className='file-list-container'>
			<h2 className='file-list-title'>File List</h2>
			<table className='file-list'>
				<thead>
					<tr className='file-list-head'>
						<th>File Name</th>
						<th>File Size</th>
						<th>Upload Time</th>
						<th>Download</th>
					</tr>
				</thead>
				<tbody>
					{fileList.map((file) => (
						<tr key={file.file_id} className='file-list-item'>
							<td className='file-name'>{file.file_name}</td>
							<td className='file-size'>{ FormatBytes(file.file_size) }</td>
							<td>{ FormatTime(file.upload_time) }</td>
							<td className='file-url'>
								<a onClick={(e) => handleDownload(e, file.file_url, file.file_name)}>
									Download
								</a>
							</td>
						</tr>
					))}
				</tbody>
			</table>
		</div>
	);
};

FileList.propTypes = {
	fileList: PropTypes.array,
};

export default FileList;