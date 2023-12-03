import React from 'react';
import { useState } from 'react';
import Input from '../Input/Input';

const RegisterForm = () => {
	const [username, setUsername] = useState('');
	const [password, setPassword] = useState('');

	const handleUsernameChange = (event) => {
		setUsername(event.target.value);
	};

	const handlePasswordChange = (event) => {
		setPassword(event.target.value);
	};

	return (
		<div>
			<form>
				<h2>Register</h2>
				<Input label='username' name='username' value={username}
					onChange={handleUsernameChange} />
				<Input label='password' name='password' value={password}
					onChange={handlePasswordChange} />
				<button type='submit'>Register</button>
			</form>
		</div>
	);
};

export default RegisterForm;