import { createDockerDesktopClient } from '@docker/extension-api-client';
let dockerClient = null;

try {
    dockerClient = createDockerDesktopClient();
} catch {

}

const DockerClient = dockerClient ? dockerClient : {
    extension: {
        vm: {
            service: {
                get: async (url) => {
                    return await fetch(`http://localhost:8080${url}`, {
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    }).then(res => res.json());
                },
                post: async (url, data) => {
                    return await fetch(`http://localhost:8080${url}`, {
                        method: 'POST',
                        body: JSON.stringify(data),
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    }).then(res => res.json());
                }
            },
        },
    },
}

export default DockerClient