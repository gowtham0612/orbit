# orbit

An open-source, self-hosted realtime event mesh written in Go, powered by Redis and WebSockets.

## Architecture

Orbit's broker architecture follows a strictly enforced `ORBIT_BROKER=redis` (V1 Engine) target.
*   **WebSocket Gateway**: Handles high-concurrency client connections. Slow consumers are forcefully disconnected preventing I/O backpressure.
*   **V1 Engine (Redis PubSub)**: Orbit features an ultra-resilient, multiplexed `PubSub` engine. Instead of a goroutine-per-channel block, Orbit streams events through a dedicated worker pool (configurable via `ORBIT_FANOUT_WORKERS=100`) heavily optimizing RAM bounds, whilst strictly preserving sequential per-channel message ordering via checksum hashing. Automatic Network-Partition reconstruction guarantees reconnection persistence.
*   **Presence Tracker**: Tracks user online status per channel using explicit events (`presence.joined`, `presence.left`) built on top of Redis Sorted Sets.

> Future Broker Adapters (V2 Redis Streams for persisted event-sourcing and V3 NATS for raw C10M scale-out speed) are architecturally viable utilizing the internal `pubsub.Engine` interface, but are explicitly unscaffolded for standard production.

## Quickstart

1.  Start the cluster using Docker Compose:
    ```bash
    docker-compose up --build
    ```
    This spins up a Redis instance and the Orbit server on port `8080`.

2.  Open `sdk/js/index.html` in your browser to interact with the event mesh.

## Example Usage (JS)

```javascript
import { Orbit } from './orbit.js';

// Connect with a token (optional auth)
const orbit = new Orbit('ws://localhost:8080/ws?token=mysecretToken');

orbit.onConnected(() => {
    console.log("Connected to Orbit!");

    // Subscribe to a channel
    orbit.subscribe('room-1', (message) => {
        console.log('Received:', message.payload);
    });

    // Publish to a channel
    orbit.publish('room-1', {
        event: 'message.created',
        payload: { text: 'Hello, World!' }
    });
});
```

## Observability

Monitor real-time metrics at `http://localhost:8080/metrics`
