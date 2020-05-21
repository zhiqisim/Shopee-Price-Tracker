import React from 'react';
import { BrowserRouter, Route, Switch, Redirect } from 'react-router-dom';
import { createMuiTheme, ThemeProvider} from '@material-ui/core/styles';

// Utils
import {AuthProvider} from './utils/Auth';
import PrivateRoute from './utils/PrivateRoute';

// Containers
import Login from './containers/Login';
import Signup from './containers/Signup';
import ItemList from './containers/ItemList';
import WatchList from './containers/WatchList';
import Price from './containers/Price';



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

    <ThemeProvider theme={theme}>
      <AuthProvider>
        <BrowserRouter>
          <div>
            <Switch>
              <PrivateRoute path="/watchlist" component={WatchList} />
              <Route path="/" exact component={ItemList} />
              <Route path="/login" exact component={Login} />
              <Route path="/signup" exact component={Signup} />
              <Route path="/price/:id" exact component={Price} />
              <Redirect from="*" to="/" />
            </Switch>
          </div>
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
