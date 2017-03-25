/**
 * @typedef {Object} txn 
 * @property {string} id
 * @property {string} state
 * @property {?string} reason
 */

(function () {
    "use strict";

    /** @type {txn} */
    var txn = JSON.parse($('#transaction-data').text());

    // This file is only loaded after the user has successfully paid but not finished vending.
    if (txn.state !== "paid") {
        console.error("This shouldn't happen? %o", txn);
        return;
    }

    /** @constant */
    var url = "/txns/" + txn.id + ".json";

    var fetchIntervalI;

    /**
     * @param {txn} newTxn 
     */
    function handleUpdate(newTxn) {
        if (newTxn.state === txn.state) {
            return
        }
        txn = newTxn;

        var newWords = txn.state;
        if (txn.state === "complete") {
            newWords = "Successfully Vended!";
        } else if (txn.state === "failed") {
            newWords = "Unable to vend: " + txn.reason;
        }
        $('#txn-state').text(newWords);

        // We're done here
        window.clearInterval(fetchIntervalI);
    }

    function fetcher() {
        $.getJSON(url).then(
            function (newTxn) {
                handleUpdate(newTxn);
            },
            function (err) {
                console.error(err);
                alert("Unable to update: " + err);
            }
        )
    }

    fetchIntervalI = setInterval(fetcher, 500);
})();
