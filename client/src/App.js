import React, { useEffect } from 'react';
import { useState } from 'react';
import './App.css';
import AuthenticationForm from './AuthenticationForm/AuthenticationForm.js';
import Dashboard from './Dashboard/Dashboard.js';

const App = () => {
	const [token, setToken] = useState(null);

	useEffect(() => {
		const token = localStorage.getItem('token');
		if (token) {
			setToken(token);
		}
	}, []);

	const handleTokenChange = (token) => {
		localStorage.setItem('token', token);
		setToken(token);
	};

	return (
		<div className="wrapper">
			{token == null ? <AuthenticationForm handleTokenChange={handleTokenChange} /> :
				<Dashboard token={token} />}
		</div>
	);
};

export default App;