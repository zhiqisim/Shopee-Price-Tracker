import React, { useContext } from 'react';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import { Button, makeStyles } from "@material-ui/core";
import { Link } from 'react-router-dom';

// import { AuthContext } from "../helpers/Auth";


const useStyles = makeStyles(theme => ({
    appBar: {
        borderBottom: `1px solid ${theme.palette.divider}`,
        background : '#FF4500',
    },
    toolbar: {
        flexWrap: 'wrap',
    },
    toolbarTitle: {
        flexGrow: 1,
    },
    link: {
        margin: theme.spacing(1, 1.5),
        textDecoration: "none",
        outlineColor: "white",
        color: "white",
    },
    logoLink: {
        textDecoration: "none",
        color: "white",
    },
}));

function HeaderButtons() {
    const classes = useStyles();
    // const { currentUser } = useContext(AuthContext);
    // if (currentUser) {
    //     return (
    //         <React.Fragment>
    //             <Link to={'/watchlist'} className={classes.link}>
    //                 <Button size="medium">
    //                     <b>POST</b>
    //                 </Button>
    //             </Link>
    //             <Link to={'/'} className={classes.link}>
    //                 <Button onClick={() => firebase.auth().signOut()} color="primary" variant="outlined" className={classes.link}>
    //                     Logout
    //                 </Button>
    //             </Link>
    //         </React.Fragment>
    //     )
    // } else {
        return (
            <Link to={'/login'} className={classes.link}>
                <Button variant="outlined" className={classes.link}>
                    <b>Login</b>
                </Button>
            </Link>
        )
    // }
}

export default function Header() {
    const classes = useStyles();
    return (
        <AppBar position="sticky" color="default" elevation={2} className={classes.appBar}>
            <Toolbar className={classes.toolbar}>
                <Typography variant="h6" color="inherit" noWrap className={classes.toolbarTitle}>
                    <Link to={'/'} className={classes.logoLink}>
                        <b>Shopee Price Tracker</b>
                    </Link>
                </Typography>
                <HeaderButtons/>
            </Toolbar>
        </AppBar>
    )

}