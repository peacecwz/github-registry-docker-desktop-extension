import * as React from 'react';
import Box from '@mui/material/Box';
import Collapse from '@mui/material/Collapse';
import IconButton from '@mui/material/IconButton';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import {Link} from "@mui/material";
import {v1} from "@docker/extension-api-client-types";
import DeleteIcon from '@mui/icons-material/Delete';
import DownloadIcon from '@mui/icons-material/Download';
import { GithubUser, Package, PackageVersion } from './types';

interface RowProps {
    readonly setState: React.Dispatch<React.SetStateAction<boolean>>,
    readonly onChange: () => void,
    readonly user: GithubUser,
    readonly pkg: Package,
    readonly client: v1.DockerDesktopClient
}

function Row(props: RowProps) {
    const {setState, user, pkg, client, onChange} = props;
    const [open, setOpen] = React.useState(false);

    const handleDownload = async (tag: PackageVersion) => {
        setState(true)

        try {
            const loginResult = await client.docker.cli.exec("login", [
                "ghcr.io",
                "--username",
                user.login,
                "--password",
                user.token
            ]);

            if (loginResult.stdout.indexOf("Login Succeeded") === -1) {
                client.desktopUI.toast.error("Cannot login to Github Registry");
                return
            }

            const pullImageResult = await client.docker.cli.exec("pull", [
                tag.image_url
            ]);

            if (pullImageResult.stdout.indexOf("Status: Downloaded newer") === -1 && pullImageResult.stdout.indexOf("Status: Image is up to date") === -1) {
                client.desktopUI.toast.error("Cannot pull from Github Registry");
                return
            }

            client.desktopUI.toast.success("Image pulled successfully");
        } catch (e) {
            client.desktopUI.toast.error("Cannot pull image from Github Registry")
        }

        setState(false)
    }

    const handleDelete = async (tag: PackageVersion) => {
        setState(true)
        try {
            await client.extension.vm?.service?.get(`/package-delete?organizationName=${tag.organizationName}&packageName=${pkg.name}&packageVersionId=${tag.id}`);
            client.desktopUI.toast.success("Deleted image successfully");
            await onChange()
        } catch (e) {
            client.desktopUI.toast.error("Cannot delete image from Github Registry")
        }
        setState(false)
    }

    return (
        <React.Fragment>
            <TableRow sx={{'& > *': {borderBottom: 'unset'}}}>
                <TableCell>
                    <IconButton
                        aria-label="expand row"
                        size="small"
                        onClick={() => setOpen(!open)}
                    >
                        {open ? <KeyboardArrowUpIcon/> : <KeyboardArrowDownIcon/>}
                    </IconButton>
                </TableCell>
                <TableCell component="th" scope="row">
                    {pkg.id}
                </TableCell>
                <TableCell component="th" scope="row">
                    <Link onClick={(e) => {
                        e.preventDefault()
                        client.host.openExternal(pkg.url)
                    }} href={pkg.url}>
                        {pkg.name}
                    </Link>
                </TableCell>
                <TableCell align="right"></TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{paddingBottom: 0, paddingTop: 0}} colSpan={6}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{margin: 1}}>
                            <Typography variant="h6" gutterBottom component="div">
                                Package Versions
                            </Typography>
                            <Table size="small" aria-label="purchases">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Id</TableCell>
                                        <TableCell>Name</TableCell>
                                        <TableCell align="right"></TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {pkg.package_versions.map((tag) => (
                                        <TableRow key={tag.id}>
                                            <TableCell component="th" scope="row">
                                                {tag.id}
                                            </TableCell>
                                            <TableCell>
                                                <Link onClick={(e) => {
                                                    e.preventDefault()
                                                    client.host.openExternal(tag.url)
                                                }} href={tag.url}>
                                                    {tag.name}
                                                </Link>
                                            </TableCell>
                                            <TableCell align="right">
                                                <IconButton onClick={async () => await handleDownload(tag)}
                                                            aria-label="download">
                                                    <DownloadIcon/>
                                                </IconButton>

                                                <IconButton onClick={async () => await handleDelete(tag)}
                                                            aria-label="delete">
                                                    <DeleteIcon/>
                                                </IconButton>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    );
}

export default function CollapsibleTable({packages, client, user, setState, onChange}) {
    return (
        <TableContainer component={Paper}>
            <Table aria-label="containers">
                <TableHead>
                    <TableRow>
                        <TableCell/>
                        <TableCell>Id</TableCell>
                        <TableCell>Name</TableCell>
                        <TableCell align="right"></TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {packages.map((pkg) => (
                        <Row key={pkg.id} onChange={onChange} setState={setState} user={user} client={client}
                             pkg={pkg}/>
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    );
}