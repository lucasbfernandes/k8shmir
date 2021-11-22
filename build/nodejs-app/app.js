const express = require('express')
const morgan = require('morgan')
const bp = require('body-parser')

const process = require('process')

process.on('SIGTERM', () => {
    console.info("Shutting down...")
    process.exit(0)
})

const app = express()

app.use(morgan('combined'))
app.use(bp.json())
app.use(bp.urlencoded({ extended: true }))
app.use((req, res, next) => {
    console.log(req.body)
    console.log(req.headers)
    next()
})

const port = 3000

let integer = 0
let lastIndex = 0

app.get('/last-index', (req, res) => {
    res.status(200).send({ index: lastIndex })
})

app.get('/integer', (req, res) => {
    lastIndex = parseInt(req.headers['log-index'])
    res.status(200).send({ value: integer })
})

app.post('/integer', (req, res) => {
    lastIndex = parseInt(req.headers['log-index'])
    const { op, value } = req.body
    switch (op) {
        case 'INC':
            integer += value
            break
        case 'DEC':
            integer -= value
            break
        default:
            res.status(500).send({ error: 'Invalid operation' });
            break
    }

    res.status(201).send()
})

app.post('/integer/reset', (req, res) => {
    lastIndex = parseInt(req.headers['log-index'])
    integer = 0
    res.status(200).send()
})

app.listen(port, () => {
    console.log(`App running on port ${port}`)
})
