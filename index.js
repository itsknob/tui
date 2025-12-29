#!/usr/bin/env node

if (process.argv.length < 3) {
  console.error("Requires Input of type Transaction")
  return 1;
}

let api = require('@actual-app/api');
require("dotenv").config();

const BANK_OF_MOM = "Bank of Mom";
const KAIDEN = '9b80c948-e1ac-499f-b8f6-4539b9ee5964';
const FINN = '0d21b510-7d09-4cbf-bd6a-a537e9a9eee1';
const AVA = 'd1ca5991-8468-46c7-9a0f-ac69562f65b9';
const accounts = {
  'kaiden': KAIDEN,
  'finn': FINN,
  'ava': AVA,
}

async function setup() {
  await api.init({
    dataDir: 'budget',
    serverURL: process.env.BUDGET_SERVER_URL,
    password: process.env.BUDGET_PASS
  })
  // console.log("Initialized app")
  await api.downloadBudget(process.env.BUDGET_SYNC_ID)
}

/**
  * @typedef Transaction
  * @property {String} account
  * @property {String} date 'YYYY-MM-DD'
  * @property {Number} amount Tranaction Amount in Cents
  * @property {String?} notes Notes for transaction
  * @property {String?} payee_name
  */

/**
  * @param {Transaction} transaction Transaction to be imported
  */
async function importTransaction(transaction) {

  // Map User Name to Account Id
  if (Object.keys(accounts).includes(transaction.account.toLowerCase())) {
    transaction.account = accounts[transaction.account.toLowerCase()]
    // console.log(transaction.account)
  }

  // Hardcode Bank of Mom
  transaction.payee_name = BANK_OF_MOM

  // Ensure format of amount is in format actual uses
  let amtString = transaction.amount.toString()
  if (amtString.includes('.')) {
    transaction.amount = parseInt(amtString.replace(".", ""))
  }

  // Confirmation
  // console.log("Final Transaction", transaction);

  // Upload
  await api.importTransactions(transaction.account, [transaction])
  // console.log("Transaction Uploaded")
}


(async (transactionString) => {
  const transaction = JSON.parse(transactionString)
  // console.log("Initial Transaction: ", transaction)
  await setup();

  await importTransaction(transaction)
  await api.shutdown()
})(process.argv[2]);

