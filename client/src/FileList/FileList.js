import React from 'react';
import './FileList.css';
import PropTypes from 'prop-types';
import FormatBytes from '../utils/formatBytes.js';

const FileList = ({ fileList }) => {
	return (
		<div className='file-list-container'>
			<h2 className='file-list-title'>File List</h2>
			<table className='file-list'>
				<thead>
					<tr className='file-list-head'>
						<th>File Name</th>
						<th>File Size</th>
						<th>Download</th>
					</tr>
				</thead>
				<tbody>
					{fileList.map((file) => (
						<tr key={file.file_id} className='file-list-item'>
							<td className='file-name'>{file.file_name}</td>
							<td className='file-size'>{ FormatBytes(file.file_size) }</td>
							<td className='file-url'>
								<a href={file.file_url}>Download</a>
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