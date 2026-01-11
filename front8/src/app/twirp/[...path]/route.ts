import { NextRequest, NextResponse } from 'next/server';
import { AuthClient } from '@/lib/api-clients';

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

    let userId = "";

    const authHeader = request.headers.get('authorization') || '';
    const authToken = authHeader.startsWith('Bearer ') 
      ? authHeader.substring(7) 
      : '';
    
    if (authToken) {
      try {
        const verifyResponse = await AuthClient.VerifyToken({ token: authToken });
        userId = verifyResponse.userId;
      } catch (error) {
        console.error('Token verification error:', error);
      }
    }

    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      'x-user-id': userId,
    };

    const body = await request.arrayBuffer();

    const response = await fetch(targetUrl, {
      method: 'POST',
      headers,
      body,
    });

    const responseBody = await response.arrayBuffer();
    
    // Forward the response with the same status code, always as JSON
    return new NextResponse(responseBody, {
      status: response.status,
      headers: {
        'Content-Type': 'application/json',
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
