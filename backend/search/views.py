import requests
from django.http import HttpResponseBadRequest
from rest_framework import status
from rest_framework.response import Response
from rest_framework.views import APIView


class SearchProducts(APIView):

    def get(self, request):

        try:
            params = self.__get_params(request.query_params)
        except ValueError as e:
            print(f'Failed to get params: {e}')
            return HttpResponseBadRequest()

        try:
            response = requests.get('http://localhost:8001/search', params=params)
            response.raise_for_status()
        except requests.exceptions.RequestException:
            return HttpResponseBadRequest()

        return Response(response.json(), status=status.HTTP_200_OK)

    @staticmethod
    def __get_params(request_params):
        params = {}
        try:
            params['phrase'] = request_params['phrase']
            params['website'] = request_params['website']
        except KeyError:
            raise ValueError('Not enough params')
        params['page'] = request_params['page'] if 'page' in request_params else 0
        return params
