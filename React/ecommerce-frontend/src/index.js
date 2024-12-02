// index.js

import React from 'react';
import { createRoot } from 'react-dom/client'; // Import createRoot
import { Provider } from 'react-redux';
import store from './store';
import App from './App';

const container = document.getElementById('root');
const root = createRoot(container); // Create a root.

root.render(
  <Provider store={store}>
    <App />
  </Provider>
);
