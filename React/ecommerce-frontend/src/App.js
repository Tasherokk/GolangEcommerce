// import React, { useEffect } from 'react';
// import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
// import { useDispatch } from 'react-redux';
// import { loadUser } from './store/actions/authActions';
// import Login from './pages/Login';
// import Register from './pages/Register';
// import HomePage from './pages/HomePage';
// import NotFound from './pages/NotFound';
// import PrivateRoute from './components/common/PrivateRoute';
// import AdminDashboard from './components/Admin/Dashboard';
// import ModeratorDashboard from './components/Moderator/Dashboard';
// import SellerDashboard from './components/Seller/Dashboard';
// import CustomerDashboard from './components/Customer/Dashboard';
// import Navbar from './components/common/Navbar';

// function App() {
    

//     return (
//         <Router>
//             <Navbar />
//             <Switch>
//                 <Route exact path="/" component={HomePage} />
//                 <Route path="/login" component={Login} />
//                 <Route path="/register" component={Register} />

//                 {/* Приватные маршруты */}
//                 <PrivateRoute path="/admin/dashboard" component={AdminDashboard} roles={['Admin']} />
//                 <PrivateRoute path="/moderator/dashboard" component={ModeratorDashboard} roles={['Moderator']} />
//                 <PrivateRoute path="/seller/dashboard" component={SellerDashboard} roles={['Seller']} />
//                 <PrivateRoute path="/customer/dashboard" component={CustomerDashboard} roles={['Customer']} />

//                 <Route component={NotFound} />
//             </Switch>
//         </Router>
//     );
// }

// export default App;


// src/App.js

import React, { useEffect } from 'react'; // Import React and useEffect
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'; // Import Router components
import { useDispatch } from 'react-redux'; // Import useDispatch from react-redux
import { loadUser } from './store/actions/authActions'; // Import loadUser action

// Import your components
import Navbar from './components/common/Navbar';
import HomePage from './pages/HomePage';
import Login from './pages/Login';
import Register from './pages/Register';
import NotFound from './pages/NotFound';

// Import your dashboard components
import AdminDashboard from './components/Admin/Dashboard';
import ModeratorDashboard from './components/Moderator/Dashboard';
import SellerDashboard from './components/Seller/Dashboard';
import CustomerDashboard from './components/Customer/Dashboard';
import ProductDetails from './components/common/ProductDetails';

// Import your PrivateRoute component
import PrivateRoute from './components/common/PrivateRoute';

function App() {
    const dispatch = useDispatch();

    useEffect(() => {
        dispatch(loadUser());
    }, [dispatch]);

    return (
        <Router>
            <Navbar />
            <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/products/:id" element={<ProductDetails />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />

                {/* Private Routes */}
                <Route
                    path="/admin/dashboard/*"
                    element={
                        <PrivateRoute roles={['Admin']}>
                            <AdminDashboard />
                        </PrivateRoute>
                    }
                />
                <Route
                    path="/moderator/dashboard/*"
                    element={
                        <PrivateRoute roles={['Moderator']}>
                            <ModeratorDashboard />
                        </PrivateRoute>
                    }
                />
                <Route
                    path="/seller/dashboard/*"
                    element={
                        <PrivateRoute roles={['Seller']}>
                            <SellerDashboard />
                        </PrivateRoute>
                    }
                />
                <Route
                    path="/customer/dashboard/*"
                    element={
                        <PrivateRoute roles={['Customer']}>
                            <CustomerDashboard />
                        </PrivateRoute>
                    }
                />
                <Route path="*" element={<NotFound />} />
            </Routes>
        </Router>
    );
}


export default App;