import { Box, Button, Drawer, List, ListItem, ListItemButton, ListItemText, TextField, Typography } from "@mui/material"
import { forwardRef, useEffect, useState } from "react"
import { useApi } from "../context/apiContext"
import { Domain } from "@concurrent-world/client/dist/types/model/core"

export const Hosts = forwardRef<HTMLDivElement>((props, ref): JSX.Element => {

    const { api } = useApi()

    const [hosts, setHosts] = useState<Domain[]>([])
    const [remoteFqdn, setRemoteFqdn] = useState('')

    const [selectedHost, setSelectedHost] = useState<Domain | null>(null)
    const [newTag, setNewTag] = useState<string>('')
    const [newScore, setNewScore] = useState<number>(0)

    useEffect(() => {
        api.getDomains().then(setHosts)
    }, [])

    return (
        <div ref={ref} {...props}>
            <Box
                width="100%"
            >
                <Box sx={{ display: 'flex', gap: '10px' }}>
                    <TextField
                        label="remote fqdn"
                        variant="outlined"
                        value={remoteFqdn}
                        sx={{ flexGrow: 1 }}
                        onChange={(e) => {
                            setRemoteFqdn(e.target.value)
                        }}
                    />
                    <Button
                        variant="contained"
                        onClick={(_) => {
                            api.sayHello(remoteFqdn)
                        }}
                    >
                        GO
                    </Button>
                </Box>
                <Typography>Hosts</Typography>
                <List
                    disablePadding
                >
                    {hosts.map((host) => (
                        <ListItem key={host.ccid}
                            disablePadding
                        >
                            <ListItemButton
                                onClick={() => {
                                    setNewTag(host.tag)
                                    setNewScore(host.score)
                                    setSelectedHost(host)
                                }}
                            >
                                <ListItemText primary={host.fqdn} secondary={`${host.ccid}`} />
                                <ListItemText>{`${host.tag}(${host.score})`}</ListItemText>
                            </ListItemButton>
                        </ListItem>
                    ))}
                </List>
            </Box>
            <Drawer
                anchor="right"
                open={selectedHost !== null}
                onClose={() => {
                    setSelectedHost(null)
                }}
            >
                <Box
                    width="50vw"
                    display="flex"
                    flexDirection="column"
                    gap={1}
                    padding={2}
                >
                    <Typography>{selectedHost?.ccid}</Typography>
                    <pre>{JSON.stringify(selectedHost, null, 2)}</pre>
                    <TextField
                        label="new tag"
                        variant="outlined"
                        value={newTag}
                        sx={{ flexGrow: 1 }}
                        onChange={(e) => {
                            setNewTag(e.target.value)
                        }}
                    />
                    <TextField
                        label="new score"
                        variant="outlined"
                        value={newScore}
                        sx={{ flexGrow: 1 }}
                        onChange={(e) => {
                            setNewScore(Number(e.target.value))
                        }}
                    />
                    <Button
                        variant="contained"
                        onClick={(_) => {
                            if (!selectedHost) return
                            api.updateDomain({
                                ...selectedHost,
                                score: newScore,
                                tag: newTag
                            })
                            setSelectedHost(null)
                        }}
                    >
                        Update
                    </Button>
                    <Button
                        variant="contained"
                        onClick={(_) => {
                            if (!selectedHost) return
                            api.deleteDomain(selectedHost.fqdn)
                            setSelectedHost(null)
                        }}
                        color="error"
                    >
                        Delete
                    </Button>
                </Box>
            </Drawer>
        </div>
    )
})

Hosts.displayName = "Hosts"

