import React from 'react';
import './App.css';
import { useState } from 'react';
import LoginForm from './LoginForm/LoginForm.js';
import RegisterForm from './RegisterForm/RegisterForm.js';

const App = () => {
	const [isLoginForm, setIsLoginForm] = useState(true);

	const toggleForm = () => {
		setIsLoginForm(!isLoginForm);
	};

	return (
		<div className="wrapper">
			<div className='authentication-form'>
				{isLoginForm ? <LoginForm /> : <RegisterForm />}
			</div>
			<button onClick={toggleForm}>
				{isLoginForm ? 'Switch to Register' : 'Switch to Login'}
			</button>
		</div>
	);
};

export default App;