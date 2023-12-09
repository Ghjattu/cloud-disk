import React from 'react';
import { useState } from 'react';
import PropTypes from 'prop-types';
import Input from '../Input/Input.js';
import AuthenticationAPI from '../api/authenticationAPI.js';

const LoginForm = ({ handleTokenChange }) => {
	const [username, setUsername] = useState('');
	const [password, setPassword] = useState('');

	const handleUsernameChange = (event) => {
		setUsername(event.target.value);
	};

	const handlePasswordChange = (event) => {
		setPassword(event.target.value);
	};

	const handleSubmit = async (event) => {
		event.preventDefault();

		try {
			const resp = await AuthenticationAPI.Login(username, password);
			handleTokenChange(resp.data.token);
		} catch (error) {
			console.log(error);
		}
	};

	return (
		<div>
			<form onSubmit={handleSubmit}>
				<h2>Login</h2>
				<Input label='username' name='username' value={username}
					onChange={handleUsernameChange} />
				<Input label='password' name='password' value={password}
					onChange={handlePasswordChange} />
				<button type='submit'>Login</button>
			</form>
		</div>
	);
};

LoginForm.propTypes = {
	handleTokenChange: PropTypes.func.isRequired,
};

export default LoginForm;