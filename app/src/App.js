import React from 'react';
import { BrowserRouter, Route, Switch, Redirect } from 'react-router-dom';
import { createMuiTheme, ThemeProvider} from '@material-ui/core/styles';
import ItemList from './containers/ItemList';
import {AuthProvider} from './helpers/Auth';
import Login from './containers/Login';
import PrivateRoute from './helpers/PrivateRoute';


const theme = createMuiTheme({
  palette: {
    primary: {
      // light: will be calculated from palette.primary.main,
      main: '#ff4400',
      // dark: will be calculated from palette.primary.main,
      // contrastText: will be calculated to contrast with palette.primary.main
    },
    secondary: {
      light: '#0066ff',
      main: '#0044ff',
      // dark: will be calculated from palette.secondary.main,
      contrastText: '#ffcc00',
    },
    // Used by `getContrastText()` to maximize the contrast between
    // the background and the text.
    contrastThreshold: 3,
    // Used by the functions below to shift a color's luminance by approximately
    // two indexes within its tonal palette.
    // E.g., shift from Red 500 to Red 300 or Red 700.
    tonalOffset: 0.2,
  },
});

function App() {
  return (

    <ThemeProvider>
      <AuthProvider>
        <BrowserRouter>
          <div>
            <Switch>
              {/* <PrivateRoute path="/watchlist" component={Post} /> */}
              <Route path="/" exact component={ItemList} />
              <Route path="/login" exact component={Login} />
              <Redirect from="*" to="/" />
            </Switch>
          </div>
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
