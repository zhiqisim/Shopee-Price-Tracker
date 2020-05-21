import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import ListItemText from '@material-ui/core/ListItemText';
import Typography from '@material-ui/core/Typography';
import API from "../utils/API";

import Header from '../components/Header';

const useStyles = makeStyles((theme) => ({
    root: {
        width: '100%',
        backgroundColor: theme.palette.background.paper,
    },
    listContainer: {
        paddingTop: theme.spacing(8),
        paddingBottom: theme.spacing(8),
    },
    card: {
        height: '100%',
        width: '100%',
        display: 'flex',
        flexDirection: 'column',
        boxShadow: theme.shadows[2],
        marginBottom: theme.spacing(4),
        marginTop: theme.spacing(4),
    },
    cardContent: {
        flexGrow: 1,
    },
    inline: {
        display: 'inline',
    },
    cardButton: {
        horizontalAlign: "right",
    },
    link: {
        margin: theme.spacing(1, 1.5),
        textDecoration: "none"
    },
}));

export default function Price(props) {
    const classes = useStyles();

    const { id } = props.match.params
    // Set default as null so we
    // know if it's still loading
    const [priceInfo, setPrice] = useState(null);

    // Initialize with listening to our
    // API. The second argument
    // with the empty array makes sure the
    // function only executes once
    useEffect(() => {
        listenForPriceInfo();
    });

    const listenForPriceInfo = () => {
        const params = {
            itemid: id,
        };
        API.get('/item/price', { params })
            .then(response => {
                console.log(response.data);
                if (response.data.message === "success") {
                    const allPrice = [];
                    const jsonData = response.data.items;
                    if (jsonData != null) {
                        jsonData.forEach((doc) => {
                            if (doc.flash_sale) {
                                doc.flash = "Yes";
                            } else {
                                doc.flash = "No";
                            }
                            allPrice.push(doc)
                        });
                    }

                    setPrice(allPrice)
                }
            })
            .catch(error => {
                console.log(error);
            });
    }

    if (!priceInfo) {
        return (
            <div>
                Please wait...
            </div>
        )
    }



    return (
        <React.Fragment>
            <Header classes={classes} />
            <List className={classes.root}>
                <Container className={classes.listContainer} maxWidth="xs">
                    {priceInfo.map(({ price_datetime, price, flash }, index) => (
                        <React.Fragment key={index}>
                            <Card className={classes.card}>
                                <CardContent className={classes.cardContent}>
                                    <ListItem alignItems="flex-start">
                                        <ListItemText
                                            primary={<React.Fragment>
                                                <Typography
                                                    component="span"
                                                    variant="body1"
                                                    className={classes.inline}
                                                    color="textPrimary"
                                                >
                                                    Price: ${price / 100000.0}

                                                </Typography>

                                            </React.Fragment>}
                                            secondary={
                                                <React.Fragment>
                                                    <br />
                                                    <Typography
                                                        component="span"
                                                        variant="body2"
                                                        className={classes.inline}
                                                        color="textPrimary"
                                                    >
                                                        {price_datetime}
                                                    </Typography>
                                                    <br />
                                                    <br />
                                                    <Typography
                                                        component="span"
                                                        variant="body2"
                                                        className={classes.inline}
                                                        color="textPrimary"
                                                    >
                                                        Flash Deal: {flash}
                                                    </Typography>

                                                </React.Fragment>
                                            }
                                        />
                                    </ListItem>
                                </CardContent>
                            </Card>

                        </React.Fragment>

                    ))}
                </Container>
            </List>
        </React.Fragment>
    );
}
