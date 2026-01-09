import { NextRequest, NextResponse } from 'next/server';

// Service name to port mapping
const SERVICE_PORTS: Record<string, number> = {
  'meme.auth.Auth': 8081,
  'meme.users.Users': 8082,
  'meme.subscriptions.Subscriptions': 8083,
  'meme.likes.Likes': 8084,
  'meme.posts.Posts': 8085,
  'meme.photos.Photos': 8086,
};

function getBackendUrl(service: string): string {
  const port = SERVICE_PORTS[service] || 8080;
  
  // In production, you might want to use environment variables
  const baseUrl = process.env.BACKEND_BASE_URL || 'http://localhost';
  return `${baseUrl}:${port}`;
}

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  try {
    const { path } = await params;
    const fullPath = path.join('/');
    
    // Parse service and method from path: "meme.auth.Auth/Login"
    const [service, method] = fullPath.split('/');
    
    if (!service || !method) {
      return NextResponse.json(
        { msg: 'Invalid path format' },
        { status: 400 }
      );
    }

    // Get the backend URL for this service
    const backendUrl = getBackendUrl(service);
    const targetUrl = `${backendUrl}/twirp/${fullPath}`;

    // Get request body as ArrayBuffer to handle both text and binary data
    const body = await request.arrayBuffer();
    
    // Get headers (forward Authorization and Content-Type)
    const contentType = request.headers.get('content-type') || 'application/json';
    const headers: HeadersInit = {
      'Content-Type': contentType,
    };
    
    const authHeader = request.headers.get('authorization');
    if (authHeader) {
      headers['Authorization'] = authHeader;
    }

    // Forward the request to the backend service
    const response = await fetch(targetUrl, {
      method: 'POST',
      headers,
      body,
    });

    // Get response body as ArrayBuffer to handle both text and binary data
    const responseBody = await response.arrayBuffer();
    const responseContentType = response.headers.get('content-type') || 'application/json';
    
    // Forward the response with the same status code and CORS headers
    return new NextResponse(responseBody, {
      status: response.status,
      headers: {
        'Content-Type': responseContentType,
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'POST, OPTIONS',
        'Access-Control-Allow-Headers': 'Content-Type, Authorization',
      },
    });
  } catch (error) {
    console.error('Proxy error:', error);
    return NextResponse.json(
      { msg: 'Internal server error' },
      { status: 500 }
    );
  }
}

// Handle OPTIONS for CORS preflight
export async function OPTIONS() {
  return new NextResponse(null, {
    status: 200,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'POST, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
    },
  });
}

