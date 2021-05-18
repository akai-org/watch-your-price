from rest_framework import status
from rest_framework.response import Response
from rest_framework.views import APIView


class SearchProducts(APIView):

    def get(self, request):
        return Response({"search": "result"}, status=status.HTTP_200_OK)
