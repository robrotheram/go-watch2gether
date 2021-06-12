import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import './index.less';

import {
  Switch,
  Route,
  Redirect,
} from 'react-router-dom';
import { Provider, connect } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';
import RoomPage from './components/RoomsPage';
import Home from './components/WelcomePage';
import reportWebVitals from './reportWebVitals';

import store, { history } from './store';

const PrivateRoute = ({ component: Component, authed, ...rest }) => {
  let redirect = '/';
  const { pathname } = rest.location;
  const lastItem = pathname.substring(pathname.lastIndexOf('/') + 1);
  if (lastItem !== 'app' || lastItem !== 'room') {
    redirect = `/?room=${lastItem}`;
  }
  return (
    <Route
      {...rest}
      render={(props) => (
        authed
          ? <Component {...props} />
          : <Redirect to={redirect} />
      )}
    />
  );
};

class Routes extends Component {
  componentDidMount() {
    console.log('==== Routes mounted!');
  }

  render() {
    return (
      <ConnectedRouter history={history}>
        <Switch>
          {/* <Route path="/app">
          <RoomPage />
        </Route> */}
          <PrivateRoute path="/app" component={RoomPage} authed={this.props.auth.auth} />
          <Route path="/">
            <Home />
          </Route>
        </Switch>
      </ConnectedRouter>
    );
  }
}
const mapStateToProps = (state) => ({ auth: state.user });
const Router = connect(mapStateToProps)(Routes);

ReactDOM.render(
  <Provider store={store}>
    <Router />
  </Provider>,
  document.getElementById('root'),
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
